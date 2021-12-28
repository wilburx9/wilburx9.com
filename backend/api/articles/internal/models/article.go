package models

import (
	"time"
)

// Article represents a single blog article
type Article struct {
	ID        string    `json:"id" firestore:"id"`
	Title     string    `json:"title" firestore:"title"`
	Thumbnail string    `json:"thumbnail" firestore:"thumbnail"`
	Url       string    `json:"url" firestore:"url"`
	PostedOn  time.Time `json:"posted_on" firestore:"posted_on"`
	UpdatedOn time.Time `json:"updated_on" firestore:"updated_on"`
	Excerpt   string    `json:"excerpt" firestore:"excerpt"`
	Source    string    `json:"source" firestore:"source"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at,serverTimestamp"`
}

// Id returns the if this Article
func (a Article) Id() string {
	return a.ID
}
