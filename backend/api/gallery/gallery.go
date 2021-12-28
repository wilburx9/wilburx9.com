package gallery

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
	"time"
)

// Handler retrieves a list of all the images sorted in descending creation date
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(internal.Database),
		HttpClient: &http.Client{},
	}

	instagram := Instagram{AccessToken: configs.Config.InstagramAccessToken, Fetch: fetch}
	unsplash := Unsplash{Username: configs.Config.UnsplashUsername, AccessKey: configs.Config.UnsplashAccessKey, Fetch: fetch}
	fetchers := [...]internal.Fetcher{instagram, unsplash}

	var images = make([]models.Image, 0)
	var updatedAts = make([]time.Time, 0)

	for _, f := range fetchers {
		var result models.ImageResult
		if err := f.GetCached(&result); err != nil {
			log.Errorf("Failed to get cached data:: %v", err)
			continue
		}
		images = append(images, result.Images...)
		updatedAts = append(updatedAts, result.UpdatedAt)
	}

	if len(images) == 0 {
		c.JSON(http.StatusInternalServerError, internal.MakeErrorResponse(images))
		return
	}

	// Sort in descending date (i.e. the most recent dates first)
	sort.Slice(images, func(i, j int) bool {
		return images[i].UploadedOn.After(images[j].UploadedOn)
	})

	c.Writer.Header().Set("Cache-Control", internal.AverageCacheControl(updatedAts))
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(images))
}