package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"time"
)

// Handler retrieves a list of all the articles
func Handler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "Work in Progress",
	})
}

// Article represents a single blog article
type Article struct {
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	Url       string    `json:"url"`
	PostedAt  time.Time `json:"posted_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Excerpt   string    `json:"excerpt"`
}

type source interface {
	fetchArticles(client common.HttpClient) []Article
}
