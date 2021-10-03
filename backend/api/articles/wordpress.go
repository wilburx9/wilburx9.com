package articles

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
)

const (
	wordpressKey = "wordpress"
)

// WordPress encapsulates fetching and caching of WordPress blog posts
type WordPress struct {
	URL string // WP V2 post URL URL e.g https://example.com/wp-json/wp/v2/posts
	internal.Fetch
}

// FetchAndCache fetches and caches WordPress articles
func (w WordPress) FetchAndCache() int {
	articles := w.fetchArticles()
	w.CacheData(internal.DbArticlesKey, wordpressKey, articles)
	return len(articles)
}

// GetCached returns cached WordPress articles
func (w WordPress) GetCached() ([]interface{}, error) {
	return w.GetCachedData(internal.DbArticlesKey, wordpressKey)
}

// fetchArticles gets articles from WordPress via HTTP
func (w WordPress) fetchArticles() []models.Article {
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
