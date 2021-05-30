package models

import (
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"time"
)

// InstaImg is container for Instagram response data
type InstaImg struct {
	Data   instaImgs `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

// InstaToken is container for Instagram token data
type InstaToken struct {
	Value       string    `json:"access_token"`
	ExpiresIn   int64     `json:"expires_in"`
	RefreshedAt time.Time `json:"refreshed_at"`
}

type instaImgs []instaImg

type instaImg struct {
	Caption   string `json:"caption"`
	MediaType string `json:"media_type"`
	ID        string `json:"id"`
	MediaURL  string `json:"media_url"`
	Timestamp string `json:"timestamp"`
	Permalink string `json:"permalink"`
}

// ToImages maps this slice of instaImg to slice of Image
func (s instaImgs) ToImages() []Image {
	var timeLayout = "2006-01-02T15:04:05-0700"
	var images = make([]Image, len(s))

	for i, e := range s {
		images[i] = Image{
			SrcThumbnail: e.MediaURL,
			Src:          e.MediaURL,
			Url:          e.Permalink,
			Caption:      e.Caption,
			UploadedAt:   internal.StringToTime(timeLayout, e.Timestamp),
			Source:       "Instagram",
		}
	}
	return images
}

// Expired returns true if this token has expired
func (t InstaToken) Expired() bool {
	var now = time.Now()
	var expireTime = t.RefreshedAt.Add(time.Second * time.Duration(t.ExpiresIn))
	return now.Equal(expireTime) || now.After(expireTime)
}

// ShouldRefresh returns true if the remaining life of the access token is less than rDuration
func (t InstaToken) ShouldRefresh(rDuration time.Duration) bool {
	var now = time.Now()
	var expireTime = t.RefreshedAt.Add(time.Second * time.Duration(t.ExpiresIn))
	diff := expireTime.Sub(now)
	return diff <= rDuration
}
