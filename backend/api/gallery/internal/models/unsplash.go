package models

import "github.com/wilburt/wilburx9.dev/backend/api/internal"

// UnsplashImgs represents a slice of unsplashImg
type UnsplashImgs []unsplashImg

type unsplashImg struct {
	CreatedAt   string `json:"created_at"`
	Color       string `json:"color"`
	Description string `json:"description"`
	User        User   `json:"User"`
	Urls        struct {
		Full  string `json:"full"`
		Small string `json:"small"`
	} `json:"urls"`
	Links struct {
		HTML string `json:"html"`
	} `json:"links"`
}

// User represents the user details of an Unsplash image
type User struct {
	Username string `json:"Username"`
	Name     string `json:"name"`
}

// ToImages maps this slice of unsplashImg to slice of Image
func (m UnsplashImgs) ToImages() []Image {
	var timeLayout = "2006-01-02T03:04:05-07:00"
	var images = make([]Image, len(m))

	for i, e := range m {
		images[i] = Image{
			SrcThumbnail: e.Urls.Small,
			Src:          e.Urls.Full,
			Url:          e.Links.HTML,
			Caption:      e.Description,
			UploadedAt:   internal.StringToTime(timeLayout, e.CreatedAt),
			Source:       "Unsplash",
			Meta: map[string]interface{}{
				"User": e.User,
			},
		}
	}
	return images
}
