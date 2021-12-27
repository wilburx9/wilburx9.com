package models

import (
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"time"
)

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

// ImageResult is data saved to the db and retrieved from it
type ImageResult struct {
	internal.Result
	Images []Image `json:"images" firestore:"images"`
}

// EmptyResponse constructs an empty ImageResult
func EmptyResponse() ImageResult {
	return ImageResult{
		Result: internal.Result{UpdatedAt: time.Now()},
		Images: []Image{},
	}
}
