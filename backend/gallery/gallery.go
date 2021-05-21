package gallery

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
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
	cacheImages()
	getCachedImages() []Image
}

func saveImages(db *badger.DB, key string, images []Image) {
	err := db.Update(func(txn *badger.Txn) error {
		buf, err := json.Marshal(images)
		if err != nil {
			return err
		}
		return txn.Set([]byte(imagesCacheKey(key)), buf)
	})
	if err != nil {
		common.LogError(fmt.Errorf("error while saving images for %v : %v", key, err))
	}
}

func getImagesFrmDb(db *badger.DB, key string) []Image {
	var images []Image
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(imagesCacheKey(key)))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &images)
		})
	})
	if err != nil {
		common.LogError(fmt.Errorf("error while getting images for %v : %v", key, err))
	}
	return images
}

func imagesCacheKey(suffix string) string {
	return common.GetCacheKey(common.StorageGallery, suffix)
}
