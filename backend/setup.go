package backend

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/articles"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"github.com/wilburt/wilburx9.dev/backend/gallery"
	"log"
	"net/http"
	"time"
)

// SetUpServer sets the Http Server. Call SetLogger before this.
func SetUpServer() *http.Server {
	gin.ForceConsoleColor()
	gin.SetMode(common.Config.Env)
	router := gin.Default()

	// Attach sentry middleware
	router.Use(sentrygin.New(sentrygin.Options{}))

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
	// Setup API route
	api := router.Group("/api")
	api.GET("/articles", articles.Handler)
	api.GET("/gallery", gallery.Handler)

	// Start Http server
	s := &http.Server{Addr: fmt.Sprintf(":%s", common.Config.ServerPort), Handler: router}
	return s
}

// SetLogger configures the custom logger
func SetLogger() error {
	log.Println(fmt.Sprintf("Sentry DSN = %s", common.Config.SentryDsn))
	return sentry.Init(sentry.ClientOptions{
		Dsn:              common.Config.SentryDsn,
		AttachStacktrace: true,
		Debug:            common.Config.Env == "debug",
		Environment:      common.Config.Env,
		TracesSampleRate: 1.0,
	})
}

// CleanUpLogger flushes buffered events
func CleanUpLogger() {
	sentry.Flush(2 * time.Second)
}
