package gallery

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
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
	fetchImages() []Image
	fetchCachedImages() []Image
	persistImages(images []Image)
}

func getImagesFrmFirestore(fsClient *firestore.Client, ctx context.Context, key string) []Image {
	cacheKey := common.GetCacheKey(common.FirestoreGallery, key)
	dSnap, err := fsClient.Collection(common.FirestoreCache).Doc(cacheKey).Get(ctx)

	if err != nil {
		if dSnap != nil && dSnap.Exists() {
			common.LogError(fmt.Errorf("error while fetching cached images for %v :: %v", key, dSnap))
		}
		return nil
	}
	var data []Image
	err = dSnap.DataTo(&data)
	if err != nil {
		common.LogError(fmt.Errorf("error while unmarshalling cached images for %v :: %v", key, dSnap))
		return nil
	}
	return data
}
