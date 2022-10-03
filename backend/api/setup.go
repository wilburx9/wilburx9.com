package api

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"github.com/wilburt/wilburx9.dev/backend/api/repos"
	"github.com/wilburt/wilburx9.dev/backend/api/update"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
)

// LoadConfig reads the configuration file and loads it into memory
func LoadConfig() error {
	return configs.LoadConfig()
}

// SetUpServer sets the Http Server. Call SetUpLogrus before this.
func SetUpServer(db database.ReadWrite) *http.Server {
	gin.ForceConsoleColor()
	gin.SetMode(configs.Config.Env)
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "It seems you are lost? Find your way buddy 😂"})
	})

	// Attach recovery middleware
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, internal.MakeErrorResponse("Something went wrong"))
	}))

	if configs.Config.IsDebug() {
		// Enable CORS support
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"http://localhost:3000"}
		corsConfig.AddAllowMethods(http.MethodGet)
		router.Use(cors.New(corsConfig))
	}

	httpClient := &http.Client{}

	update.SetUp(httpClient, db)

	// Setup API route.
	api := router.Group("/api")
	api.GET("/articles", func(c *gin.Context) { articles.Handler(c, db) })
	api.GET("/gallery", func(c *gin.Context) { gallery.Handler(c, db) })
	api.GET("/repos", func(c *gin.Context) { repos.Handler(c, db) })

	auth := api.Group("/protected")
	auth.Use(internal.AuthMiddleware())
	auth.POST("/cache", func(c *gin.Context) { update.Handler(c, db, httpClient) })

	// Start Http server
	s := &http.Server{Addr: fmt.Sprintf(":%s", configs.Config.Port), Handler: router}
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

// SetUpDatabase sets up Firebase Firestore in release  and a local db in debug
func SetUpDatabase() database.ReadWrite {
	if configs.Config.IsRelease() {
		ctx := context.Background()
		projectId := configs.Config.GcpProjectId
		client, err := firestore.NewClient(ctx, projectId)
		if err != nil {
			log.Fatalf("Failed to create Firestore cleint: %v", err)
		}
		return &database.FirebaseFirestore{
			Client: client,
			Ctx:    ctx,
		}
	} else {
		return &database.LocalDatabase{}
	}
}
