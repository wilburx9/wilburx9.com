package api

import (
	"bytes"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles"
	"github.com/wilburt/wilburx9.dev/backend/api/contact"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/repos"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var config = &configs.Config

// LoadConfig reads the configuration file and loads it into memory
func LoadConfig() error {
	return configs.LoadConfig("../configs")
}

// SetUpServer sets the Http Server. Call SetUpLogrus before this.
func SetUpServer(db *badger.DB) *http.Server {
	gin.ForceConsoleColor()
	gin.SetMode(config.Env)
	router := gin.Default()

	// Attach sentry middleware
	router.Use(sentrygin.New(sentrygin.Options{}))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "It seems you are lost? Find your way buddy 😂"})
	})

	// Attach API middleware
	router.Use(apiMiddleware(db))
	router.Use(static.Serve("/", static.LocalFile("../../frontend/build", true)))
	// Setup API route
	api := router.Group("/api")
	api.GET("/articles", articles.Handler)
	api.GET("/gallery", gallery.Handler)
	api.GET("/repos", repos.Handler)
	api.POST("/contact", contact.Handler)

	// Start Http server
	s := &http.Server{Addr: fmt.Sprintf(":%s", config.Port), Handler: router}
	return s
}

// SetUpLogrus configures the Logrus
func SetUpLogrus() {
	// Setup Logrus
	log.SetLevel(log.TraceLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		PadLevelText:  true,
	})
}

// SetUpSentry configures Sentry and attaches a Logrus hook
func SetUpSentry() error {
	var hook = internal.NewSentryLogrusHook([]log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.TraceLevel,
	})
	// Setup Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDsn,
		AttachStacktrace: true,
		Debug:            config.IsDebug(),
		Environment:      config.Env,
		TracesSampleRate: 1.0,
	})
	log.AddHook(&hook)
	return err
}

// ApiMiddleware adds custom params to request contexts
func apiMiddleware(db *badger.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(internal.Db, db)
		c.Next()
	}
}

// ScheduleFetchAddCache schedules fetching and caching of data from fetchers
func ScheduleFetchAddCache(db *badger.DB) {
	s := gocron.NewScheduler(time.UTC)
	s.Every(3).Days().Do(func(db *badger.DB) {
		fetchAndCache(db)
	}, db)
	s.StartAsync()
}

type result struct {
	fetcher string
	size    int
}

// fetchAndCache iteratively calls fetchAndCache all all fetchers
func fetchAndCache(db *badger.DB) {
	var startTime = time.Now()
	var config = &configs.Config
	fetcher := internal.Fetch{
		Db:         db,
		HttpClient: &http.Client{},
	}

	instagram := gallery.Instagram{AccessToken: config.InstagramAccessToken, Fetch: fetcher}
	unsplash := gallery.Unsplash{Username: config.UnsplashUsername, AccessKey: config.UnsplashAccessKey, Fetch: fetcher}
	medium := articles.Medium{Name: config.MediumUsername, Fetch: fetcher}
	wordpress := articles.Wordpress{URL: config.WPUrl, Fetch: fetcher}
	github := repos.Github{Auth: config.GithubToken, Username: config.UnsplashUsername, Fetch: fetcher}

	fetchers := [...]internal.Fetcher{instagram, unsplash, medium, wordpress, github}

	results := make(chan result, len(fetchers))
	var wg sync.WaitGroup

	for _, f := range fetchers {
		wg.Add(1)
		go fetchAndCacheFetcher(&wg, f, results)
	}
	wg.Wait()
	close(results)

	buffer := &bytes.Buffer{}
	for r := range results {
		buffer.WriteString(fmt.Sprintf("	%v: %d\n", r.fetcher, r.size))
	}
	var message = `
	==================== Cache Result ====================
%v	-------------------- %v ---------------------
	==================== Cache Result ====================
	`
	log.Tracef(message, buffer.String(), time.Since(startTime))
	db.RunValueLogGC(0.7)
}

func fetchAndCacheFetcher(wg *sync.WaitGroup, fetcher internal.Fetcher, out chan<- result) {
	defer wg.Done()
	size := fetcher.FetchAndCache()
	out <- result{fetcher: reflect.TypeOf(fetcher).Name(), size: size}
}
