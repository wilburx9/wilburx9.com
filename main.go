package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/articles"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"github.com/wilburt/wilburx9.dev/backend/gallery"
	"log"
	"net/http"
)

func main() {
	// Attempt to load config
	if err := common.LoadConfig("./"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}
	log.Printf("Config = %+v\\n", common.Config)
	gin.ForceConsoleColor()
	gin.SetMode(common.Config.Env)
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))

	// Setup API route
	api := router.Group("/api")
	api.GET("/articles", articles.Handler)
	api.GET("/gallery", gallery.Handler)

	s := &http.Server{Addr: fmt.Sprintf(":%s", common.Config.Port), Handler: router}
	s.ListenAndServe()
}
