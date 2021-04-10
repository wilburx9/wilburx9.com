package gallery

import (
	"github.com/wilburt/wilburx9.dev/backend/common"
	"time"
)

type Image struct {
	Thumbnail  string                   `json:"thumbnail"`
	Url        string                   `json:"url"`
	Caption    string                   `json:"caption"`
	UploadedAt time.Time                `json:"uploaded_at"`
	Source     string                   `json:"source"`
	Meta       map[string]interface{} `json:"meta"`
}

type source interface {
	fetchImages(client common.HttpClient) []Image
}
