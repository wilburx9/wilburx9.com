package backend

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/articles"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"github.com/wilburt/wilburx9.dev/backend/gallery"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"time"
)

// SetUpServer sets the Http Server. Call SetLogger before this.
func SetUpServer(dbCtx context.Context, fsClient *firestore.Client) *http.Server {
	gin.ForceConsoleColor()
	gin.SetMode(common.Config.Env)
	router := gin.Default()

	// Attach sentry middleware
	router.Use(sentrygin.New(sentrygin.Options{}))

	// Attach API middleware
	router.Use(common.ApiMiddleware(dbCtx, fsClient))

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

// CreateFireStoreClient configures and returns a pointer for Firestore client
func CreateFireStoreClient(ctx context.Context) *firestore.Client {
	var options = option.WithCredentialsJSON([]byte(common.Config.GcpSaKey))
	client, err := firestore.NewClient(ctx, common.Config.GcpProjectId, options)
	if err != nil {
		common.LogError(err)
	}
	return client
}
