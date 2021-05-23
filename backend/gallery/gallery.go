package gallery

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"sort"
	"time"
)

// Handler retrieves a list of all the images
func Handler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	fetcher := common.Fetcher{
		Db:         c.MustGet(common.Db).(*badger.DB),
		HttpClient: &http.Client{},
	}

	instagram := Instagram{AccessToken: common.Config.InstagramAccessToken, Fetcher: fetcher}
	// unsplash := Unsplash{Username: common.Config.UnsplashUsername, AccessKey: common.Config.UnsplashAccessKey, Fetcher: fetcher}
	sources := [...]common.Source{instagram}

	var allImages []Image
	for _, source := range sources {
		var images []Image
		bytes, _ := source.GetCached()
		json.Unmarshal(bytes, &images)
		allImages = append(allImages, images...)
	}

	// Sort in descending date (i.e the most recent dates first)
	sort.Slice(allImages, func(i, j int) bool {
		return allImages[i].UploadedAt.After(allImages[j].UploadedAt)
	})
	c.JSON(http.StatusOK, common.MakeSuccessResponse(allImages))
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

func getCacheKey(suffix string) string {
	return common.GetCacheKey(common.StorageGallery, suffix)
}
