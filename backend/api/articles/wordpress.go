package articles

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
	"regexp"
	"strings"
)

const (
	wordpressKey = "Wordpress"
)

// Wordpress encapsulates fetching and caching of wordpress blog posts
type Wordpress struct {
	URL string // WP V2 post URL URL e.g https://example.com/wp-json/wp/v2/posts
	internal.Fetch
}

// FetchAndCache fetches and caches wordpress articles
func (w Wordpress) FetchAndCache() {
	articles := w.FetchArticles()
	buf, _ := json.Marshal(articles)
	w.CacheData(getCacheKey(wordpressKey), buf)
}

// GetCached returns cached Wordpress articles
func (w Wordpress) GetCached() ([]byte, error) {
	return w.GetCachedData(getCacheKey(wordpressKey))
}

// FetchArticles gets articles from Wordpress via HTTP
func (w Wordpress) FetchArticles() []Article {
	req, err := http.NewRequest(http.MethodGet, w.URL, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return nil
	}

	res, err := w.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil
	}
	defer res.Body.Close()

	var posts posts
	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
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
			PostedAt:  internal.StringToTime(timeLayout, e.Date),
			UpdatedAt: internal.StringToTime(timeLayout, e.Date),
			Excerpt:   getWpExcept(e.Excerpt.Rendered),
		}
	}
	return articles
}

// Remove Html tag and leading & trailing spaces from the except
func getWpExcept(s string) string {
	var rt = regexp.MustCompile(`<[^>]*>`)   // Tags regex
	var noTags = rt.ReplaceAllString(s, " ") // Remove tags

	var rs = regexp.MustCompile(`/\\s{2,}`)        // Double spaces regex
	var noSpaces = rs.ReplaceAllString(noTags, "") // Remove double spaces

	return strings.TrimSpace(noSpaces)
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
