package models

import (
	"time"
)

// Image is a container for each object returned by Handler
type Image struct {
	ID         string                 `json:"id" firestore:"id"`
	Thumbnail  string                 `json:"thumbnail" firestore:"thumbnail"`
	Page       string                 `json:"page" firestore:"page"`
	Url        string                 `json:"url" firestore:"url"`
	Caption    string                 `json:"caption" firestore:"caption"`
	UploadedOn time.Time              `json:"uploaded_on" firestore:"uploaded_on"`
	Source     string                 `json:"source" firestore:"source"`
	Meta       map[string]interface{} `json:"meta" firestore:"meta"`
	UpdatedAt  time.Time              `json:"updated_at" firestore:"updated_at,serverTimestamp"`
}

// Id returns the if this Image
func (i Image) Id() string {
	return i.ID
}
