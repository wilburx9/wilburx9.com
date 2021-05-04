package gallery

import (
	"github.com/gin-gonic/gin"
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

// Image is a container for each object returned by Handler
type Image struct {
	SrcThumbnail string                 `json:"src_thumbnail"`
	Url          string                 `json:"url"`
	Src          string                 `json:"src"`
	Caption      string                 `json:"caption"`
	UploadedAt   time.Time              `json:"uploaded_at"`
	Source       string                 `json:"source"`
	Meta         map[string]interface{} `json:"meta"`
}

type source interface {
	fetchImages() []Image
}
