package gallery

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
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

	var allImages = make([]models.Image, 0)
	for _, f := range fetchers {
		var result []models.Image
		if err := f.GetCached(&result); err != nil {
			log.Errorf("Failed to get cached data:: %v", err)
			continue
		}
		allImages = append(allImages, result...)
	}

	if len(allImages) == 0 {
		c.JSON(http.StatusInternalServerError, internal.MakeErrorResponse(allImages))
		return
	}

	// Sort in descending date (i.e. the most recent dates first)
	sort.Slice(allImages, func(i, j int) bool {
		return allImages[i].UploadedAt.After(allImages[j].UploadedAt)
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allImages))
}