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
	result := w.fetchArticles()
	w.Db.Persist(internal.DbArticlesKey, wordpressKey, result)
	return len(result.Articles)
}

// GetCached returns cached WordPress articles
func (w WordPress) GetCached(result interface{}) error {
	return w.Db.Retrieve(internal.DbArticlesKey, wordpressKey, result)
}

// fetchArticles gets articles from WordPress via HTTP
func (w WordPress) fetchArticles() models.ArticleResult {
	req, err := http.NewRequest(http.MethodGet, w.URL, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return models.EmptyResponse()
	}

	res, err := w.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return models.EmptyResponse()
	}
	defer res.Body.Close()

	var posts models.WpPosts
	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return models.EmptyResponse()
	}
	return posts.ToResult()
}
