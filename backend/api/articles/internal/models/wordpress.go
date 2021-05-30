package models

import (
	"fmt"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"regexp"
	"strings"
)

// WpPost is a container for Wordpress post
type WpPost struct {
	Date     string  `json:"date"`
	Modified string  `json:"modified"`
	Link     string  `json:"link"`
	Title    content `json:"title"`
	Excerpt  content `json:"excerpt"`
	Meta     meta    `json:"meta"`
}
type meta struct {
	Thumbnail string `json:"twitter-card-image"`
}

type content struct {
	Rendered string `json:"rendered"`
}

// WpPosts is an slice of WpPost
type WpPosts []WpPost

// ToArticles maps this slice of WpPost to a slice of Article
func (p WpPosts) ToArticles() []Article {
	var timeLayout = "2006-01-02T15:04:05"
	var articles = make([]Article, len(p))

	for i, e := range p {
		articles[i] = Article{
			Title:     e.Title.Rendered,
			Thumbnail: e.Meta.Thumbnail,
			Url:       e.Link,
			PostedAt:  internal.StringToTime(timeLayout, e.Date),
			UpdatedAt: internal.StringToTime(timeLayout, e.Date),
			Excerpt:   fmt.Sprintf("%v..", getWpExcept(e.Excerpt.Rendered)),
		}
	}
	return articles
}

// Remove Html tag, leading and trailing spaces from the excerpt
func getWpExcept(s string) string {
	var rt = regexp.MustCompile(`<[^>]*>`)   // Tags regex
	var noTags = rt.ReplaceAllString(s, " ") // Remove tags

	var rs = regexp.MustCompile(`/\\s{2,}`)        // Double spaces regex
	var noSpaces = rs.ReplaceAllString(noTags, "") // Remove double spaces

	return strings.TrimSpace(noSpaces)
}
