package models

import (
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
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
	ID          string    `json:"id" firestore:"id"`
	Value       string    `json:"access_token" firestore:"access_token"`
	ExpiresIn   int64     `json:"expires_in" firestore:"expires_in"`
	RefreshedAt time.Time `json:"refreshed_at" firestore:"refreshed_at"`
}

// NewInstaToken returns a pointer to an InstaToken that's empty except for the id
func NewInstaToken(id string) *InstaToken {
	return &InstaToken{ID: id}
}

// Id returns the if this token
func (t InstaToken) Id() string {
	return t.ID
}

type instaImgs []instaImg

type instaImg struct {
	Caption      string `json:"caption"`
	MediaType    string `json:"media_type"`
	ID           string `json:"id"`
	MediaURL     string `json:"media_url"`
	Timestamp    string `json:"timestamp"`
	Permalink    string `json:"permalink"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

func (i instaImg) thumbnail() string {
	if i.ThumbnailUrl == "" {
		return i.MediaURL
	}
	return i.ThumbnailUrl
}

// ToImages maps this slice of instaImg to slice of Image
func (s instaImgs) ToImages(source string) []database.Model {
	var timeLayout = "2006-01-02T15:04:05-0700"
	var images = make([]database.Model, len(s))

	for i, e := range s {

		images[i] = Image{
			ID:         internal.MakeId(source, e.ID),
			Thumbnail:  e.thumbnail(),
			Page:       e.Permalink,
			Url:        e.MediaURL,
			Caption:    e.Caption,
			UploadedOn: internal.StringToTime(timeLayout, e.Timestamp),
			Source:     source,
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

// IsAboutToExpire returns true if the remaining life of the access token is less than r
func (t InstaToken) IsAboutToExpire(r time.Duration) bool {
	var now = time.Now()
	var expireTime = t.RefreshedAt.Add(time.Second * time.Duration(t.ExpiresIn))
	diff := expireTime.Sub(now)
	return diff <= r
}
