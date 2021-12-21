package models

import (
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"time"
)

// Article represents a single blog article
type Article struct {
	Title     string    `json:"title" firestore:"title" mapstructure:"title"`
	Thumbnail string    `json:"thumbnail" firestore:"thumbnail" mapstructure:"thumbnail"`
	Url       string    `json:"url" firestore:"url" mapstructure:"url"`
	PostedAt  time.Time `json:"posted_at" firestore:"posted_at" mapstructure:"posted_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at" mapstructure:"updated_at"`
	Excerpt   string    `json:"excerpt" firestore:"excerpt" mapstructure:"excerpt"`
}

// ArticleResult is data saved to the db and retrieved from it
type ArticleResult struct {
	internal.Result
	Articles []Article `json:"articles" firestore:"articles"`
}

// EmptyResponse constructs an empty ArticleResult
func EmptyResponse() ArticleResult {
	return ArticleResult{
		Result:   internal.Result{UpdatedAt: time.Now()},
		Articles: []Article{},
	}
}
