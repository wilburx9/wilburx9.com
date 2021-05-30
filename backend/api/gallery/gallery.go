package gallery

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
)

// Handler retrieves a list of all the images sorted in descending creation date
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(*badger.DB),
		HttpClient: &http.Client{},
	}

	instagram := Instagram{AccessToken: configs.Config.InstagramAccessToken, Fetch: fetch}
	unsplash := Unsplash{Username: configs.Config.UnsplashUsername, AccessKey: configs.Config.UnsplashAccessKey, Fetch: fetch}
	fetchers := [...]internal.Fetcher{instagram, unsplash}

	var allImages = make([]models.Image, 0)
	for _, f := range fetchers {
		var images []models.Image
		bytes, _ := f.GetCached()
		json.Unmarshal(bytes, &images)
		allImages = append(allImages, images...)
	}

	// Sort in descending date (i.e the most recent dates first)
	sort.Slice(allImages, func(i, j int) bool {
		return allImages[i].UploadedAt.After(allImages[j].UploadedAt)
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allImages))
}

func getCacheKey(suffix string) string {
	return internal.GetCacheKey(internal.DbGalleryKey, suffix)
}
