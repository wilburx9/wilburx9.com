package articles

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
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
	articles := w.fetchArticles()
	buf, _ := json.Marshal(articles)
	w.CacheData(getCacheKey(wordpressKey), buf)
}

// GetCached returns cached Wordpress articles
func (w Wordpress) GetCached() ([]byte, error) {
	return w.GetCachedData(getCacheKey(wordpressKey))
}

// fetchArticles gets articles from Wordpress via HTTP
func (w Wordpress) fetchArticles() []models.Article {
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

	var posts models.WpPosts
	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil
	}
	return posts.ToArticles()
}

