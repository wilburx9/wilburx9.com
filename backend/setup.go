package backend

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/articles"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"github.com/wilburt/wilburx9.dev/backend/gallery"
	"net/http"
	"time"
)

var config = common.Config

// SetUpServer sets the Http Server. Call SetLogger before this.
func SetUpServer(db *badger.DB) *http.Server {
	gin.ForceConsoleColor()
	gin.SetMode(config.Env)
	router := gin.Default()

	// Attach sentry middleware
	router.Use(sentrygin.New(sentrygin.Options{}))

	// Attach API middleware
	router.Use(common.ApiMiddleware(db))

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
	// Setup API route
	api := router.Group("/api")
	api.GET("/articles", articles.Handler)
	api.GET("/gallery", gallery.Handler)

	// Start Http server
	s := &http.Server{Addr: fmt.Sprintf(":%s", config.Port), Handler: router}
	return s
}

// SetLogger configures the custom logger
func SetLogger() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDsn,
		AttachStacktrace: true,
		Debug:            config.Env == "debug",
		Environment:      config.Env,
		TracesSampleRate: 1.0,
	})
}

// CleanUpLogger flushes buffered events
func CleanUpLogger() {
	sentry.Flush(2 * time.Second)
}

// CacheDataSources iteratively calls FetchAndCache all all data sources
func CacheDataSources(db *badger.DB) {
	fetcher := common.Fetcher{
		Db:         db,
		HttpClient: &http.Client{},
	}

	instagram := gallery.Instagram{AccessToken: config.InstagramAccessToken, Fetcher: fetcher}
	unsplash := gallery.Unsplash{Username: config.UnsplashUsername, AccessKey: config.UnsplashAccessKey, Fetcher: fetcher}
	medium := articles.Medium{Name: config.MediumUsername, Fetcher: fetcher}
	wordpress := articles.Wordpress{URL: config.WPUrl, Fetcher: fetcher}

	sources := [...]common.Source{instagram, unsplash, medium, wordpress}
	for _, source := range sources {
		go source.FetchAndCache()
	}
	db.RunValueLogGC(0.7)
}
