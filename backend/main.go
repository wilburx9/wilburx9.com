package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/articles"
	"github.com/wilburt/wilburx9.dev/backend/gallery"
	"net/http"
)

func main() {
	gin.ForceConsoleColor()

	router := gin.Default()
	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("../frontend/build", true)))

	// Setup API route
	api := router.Group("/api")
	api.GET("/articles", articles.Handler)
	api.GET("/gallery", gallery.Handler)

	var port = "8080" // TODO: Load dynamically
	s := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: router}
	s.ListenAndServe()
}
