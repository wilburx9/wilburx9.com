package articles

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"net/http"
)

const (
	wordpressKey = "wordpress"
)

// WordPress encapsulates fetching and caching of WordPress blog posts
type WordPress struct {
	URL string // WP V2 post URL URL e.g https://example.com/wp-json/wp/v2/posts
	Db         database.ReadWrite
	HttpClient internal.HttpClient
}

// Cache fetches and caches WordPress articles
func (w WordPress) Cache() (int, error) {
	result, err := w.fetchArticles()
	if err != nil {
		return 0, err
	}

	return len(result), w.Db.Write(internal.DbArticlesKey, result...)
}

// fetchArticles gets articles from WordPress via HTTP
func (w WordPress) fetchArticles() ([]database.Model, error) {
	req, _ := http.NewRequest(http.MethodGet, w.URL, nil)
	res, err := w.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil, err
	}
	defer res.Body.Close()

	var posts models.WpPosts
	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil, err
	}
	return posts.ToResult(wordpressKey), nil
}
