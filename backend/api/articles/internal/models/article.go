package models

import "time"

// Article represents a single blog article
type Article struct {
	Title     string    `json:"title" firestore:"title" mapstructure:"title"`
	Thumbnail string    `json:"thumbnail" firestore:"thumbnail" mapstructure:"thumbnail"`
	Url       string    `json:"url" firestore:"url" mapstructure:"url"`
	PostedAt  time.Time `json:"posted_at" firestore:"posted_at" mapstructure:"posted_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at" mapstructure:"updated_at"`
	Excerpt   string    `json:"excerpt" firestore:"excerpt" mapstructure:"excerpt"`
}
