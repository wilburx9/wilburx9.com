package api

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
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
	return configs.LoadConfig()
}

// SetUpServer sets the Http Server. Call SetUpLogrus before this.
func SetUpServer(db internal.Database) *http.Server {
	gin.ForceConsoleColor()
	gin.SetMode(config.Env)
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "It seems you are lost? Find your way buddy 😂"})
	})

	// Attach API middleware
	router.Use(apiMiddleware(db))

	// Attach recovery middleware
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, internal.MakeErrorResponse("Something went wrong"))
	}))

	if config.IsDebug() {
		// Enable CORS support
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"http://localhost:3000"}
		corsConfig.AddAllowMethods(http.MethodGet)
		router.Use(cors.New(corsConfig))
	}

	httpClient := &http.Client{}

	// Setup API route.
	api := router.Group("/api")
	api.GET("/articles", articles.Handler)
	api.GET("/gallery", gallery.Handler)
	api.GET("/repos", repos.Handler)
	api.POST("/contact", func(c *gin.Context) { contact.Handler(c, httpClient) })

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

// ApiMiddleware adds custom params to request contexts
func apiMiddleware(db internal.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(internal.Db, db)
		c.Next()
	}
}

// ScheduleCache schedules caching of data from cachers
func ScheduleCache(db internal.Database) {
	if config.IsDebug() {
		return
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Weeks().Do(func(db internal.Database) {
		cacheData(db)
	}, db)
	s.StartAsync()
}

// SetUpDatabase sets up Firebase Firestore in release  and a local db in debug
func SetUpDatabase() internal.Database {
	if configs.Config.IsRelease() {
		ctx := context.Background()
		projectId := configs.Config.GcpProjectId
		client, err := firestore.NewClient(ctx, projectId)
		if err != nil {
			log.Fatalf("Failed to create Firestore cleint: %v", err)
		}
		return &internal.FirebaseFirestore{
			Client: client,
			Ctx:    ctx,
		}
	} else {
		return &internal.LocalDatabase{}
	}
}

func cacheData(db internal.Database) {
	var startTime = time.Now()
	var config = &configs.Config
	fetcher := internal.Fetch{
		Db:         db,
		HttpClient: &http.Client{},
	}

	instagram := gallery.Instagram{AccessToken: config.InstagramAccessToken, Fetch: fetcher}
	unsplash := gallery.Unsplash{Username: config.UnsplashUsername, AccessKey: config.UnsplashAccessKey, Fetch: fetcher}
	medium := articles.Medium{Name: config.MediumUsername, Fetch: fetcher}
	wordpress := articles.WordPress{URL: config.WPUrl, Fetch: fetcher}
	github := repos.GitHub{Auth: config.GithubToken, Username: config.UnsplashUsername, Fetch: fetcher}

	fetchers := [...]internal.Cacher{instagram, unsplash, medium, wordpress, github}

	results := make(chan result, len(fetchers))
	var wg sync.WaitGroup

	for _, f := range fetchers {
		wg.Add(1)
		go cacheWith(&wg, f, results)
	}
	wg.Wait()
	close(results)

	buffer := &bytes.Buffer{}
	for  r := range results {
		buffer.WriteString(fmt.Sprintf("	%v: %d\n", r.cacher, r.size))
	}
	var message = `
	==================== Cache Result ====================
%v	-------------------- %v ---------------------
	==================== Cache Result ====================
	`
	log.Tracef(message, buffer.String(), time.Since(startTime))
}

func cacheWith(wg *sync.WaitGroup, cacher internal.Cacher, out chan<- result) {
	defer wg.Done()
	size := cacher.Cache()
	out <- result{cacher: reflect.TypeOf(cacher).Name(), size: size}
}


type result struct {
	cacher string
	size   int
}
