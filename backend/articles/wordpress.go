package articles

import (
	"encoding/json"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
)

type wordpress struct {
	url string // WP V2 post url URL e.g https://example.com/wp-json/wp/v2/posts
}

func (w wordpress) fetchArticles(client common.HttpClient) []Article {
	req, err := http.NewRequest(http.MethodGet, w.url, nil)
	if err != nil {
		common.LogError(err)
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		common.LogError(err)
		return nil
	}
	defer res.Body.Close()

	var posts posts
	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		common.LogError(err)
		return nil
	}
	return posts.toArticles()
}

func (p posts) toArticles() []Article {
	var timeLayout = "2006-01-02T15:04:05"
	var articles = make([]Article, len(p))
	for i, e := range p {
		articles[i] = Article{
			Title:     e.Title.Rendered,
			Thumbnail: e.Meta.Thumbnail,
			Url:       e.Link,
			PostedAt:  common.StringToTime(timeLayout, e.Date),
			UpdatedAt: common.StringToTime(timeLayout, e.Date),
			Excerpt:   e.Excerpt.Rendered,
		}
	}
	return articles
}

type posts []post

type post struct {
	Date     string  `json:"date"`
	Modified string  `json:"modified"`
	Link     string  `json:"link"`
	Title    content `json:"title"`
	Excerpt  content `json:"excerpt"`
	Meta     struct {
		Thumbnail string `json:"twitter-card-image"`
	} `json:"meta"`
}
type content struct {
	Rendered string `json:"rendered"`
}
