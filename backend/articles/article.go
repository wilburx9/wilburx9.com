package articles

import (
	"net/http"
	"time"
)

// Article represents a single blog article
type Article struct {
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	Url       string    `json:"url"`
	PostedAt  time.Time `json:"posted_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Excerpt   string    `json:"excerpt"`
}

type source interface {
	fetchArticles(client *http.Client) []Article
}