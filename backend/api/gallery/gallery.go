package gallery

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
)

// Handler retrieves a list of all the images sorted in descending creation date
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(*firestore.Client),
		Ctx:        c,
		HttpClient: &http.Client{},
	}

	instagram := Instagram{AccessToken: configs.Config.InstagramAccessToken, Fetch: fetch}
	unsplash := Unsplash{Username: configs.Config.UnsplashUsername, AccessKey: configs.Config.UnsplashAccessKey, Fetch: fetch}
	fetchers := [...]internal.Fetcher{instagram, unsplash}

	var allImages = make([]models.Image, 0)
	for _, f := range fetchers {
		maps, _ := f.GetCached()
		allImages = append(allImages, mapSliceToImages(maps)...)
	}

	// Sort in descending date (i.e. the most recent dates first)
	sort.Slice(allImages, func(i, j int) bool {
		return allImages[i].UploadedAt.After(allImages[j].UploadedAt)
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allImages))
}

func mapSliceToImages(maps []interface{}) []models.Image {
	var images = make([]models.Image, len(maps))
	for i, v := range maps {
		var image models.Image
		err := mapstructure.Decode(v, &image)
		if err == nil {
			images[i] = image
		}
	}
	return images
}
