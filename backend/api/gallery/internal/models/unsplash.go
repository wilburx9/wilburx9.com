package models

import (
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
)

// UnsplashImgs represents a slice of unsplashImg
type UnsplashImgs []unsplashImg

type unsplashImg struct {
	ID          string `json:"id"`
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
	Username string `json:"username"`
	Name     string `json:"name"`
}

// ToImages maps this slice of unsplashImg to slice of Image
func (m UnsplashImgs) ToImages(source string) []database.Model {
	var timeLayout = "2006-01-02T03:04:05-07:00"
	var images = make([]database.Model, len(m))

	for i, e := range m {
		images[i] = Image{
			ID:         internal.MakeId(source, e.ID),
			Thumbnail:  e.Urls.Small,
			Page:       e.Links.HTML,
			Url:        e.Urls.Full,
			Caption:    e.Description,
			UploadedOn: internal.StringToTime(timeLayout, e.CreatedAt),
			Source:     source,
			Meta: map[string]interface{}{
				"user": e.User,
			},
		}
	}
	return images
}
