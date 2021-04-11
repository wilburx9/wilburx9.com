package gallery

import (
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"time"
)

// Handler retrieves a list of all the images
func Handler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "Work in Progress",
	})
}

type Image struct {
	Thumbnail  string                   `json:"thumbnail"`
	Url        string                   `json:"url"`
	Caption    string                   `json:"caption"`
	UploadedAt time.Time                `json:"uploaded_at"`
	Source     string                   `json:"source"`
	Meta       map[string]interface{} `json:"meta"`
}

type source interface {
	fetchImages(client common.HttpClient) []Image
}
