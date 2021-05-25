package gallery

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
	"sort"
	"time"
)

// Handler retrieves a list of all the images sorted in descending creation date
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(*badger.DB),
		HttpClient: &http.Client{},
	}

	instagram := Instagram{AccessToken: internal.Config.InstagramAccessToken, Fetch: fetch}
	unsplash := Unsplash{Username: internal.Config.UnsplashUsername, AccessKey: internal.Config.UnsplashAccessKey, Fetch: fetch}
	fetchers := [...]internal.Fetcher{instagram, unsplash}

	var allImages = make([]Image, 0)
	for _, f := range fetchers {
		var images []Image
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
	return internal.GetCacheKey(internal.DbGalleryKey, suffix)
}
