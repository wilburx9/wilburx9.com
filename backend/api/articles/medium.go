package articles

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
)

const (
	mediumKey = "medium"
)

// Medium encapsulates the fetching and caching of medium articles
type Medium struct {
	Name string // should be Medium username (e.g "@Wilburx9") or publication (e.g. flutter-community)
	internal.Fetch
}

// FetchAndCache fetches and caches all Medium Articles
func (m Medium) FetchAndCache() int {
	result := m.fetchArticles()
	m.Db.Persist(internal.DbArticlesKey, mediumKey, result)
	return len(result.Articles)
}

// GetCached returns cached Medium articles
func (m Medium) GetCached(result interface{}) error {
	return m.Db.Retrieve(internal.DbArticlesKey, mediumKey, result)
}

// fetchArticles fetches articles via HTTP
func (m Medium) fetchArticles() models.ArticleResult {
	url := fmt.Sprintf("https://medium.com/feed/%s", m.Name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return models.EmptyResponse()
	}

	res, err := m.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return models.EmptyResponse()
	}
	defer res.Body.Close()

	var rss models.Rss
	err = xml.NewDecoder(res.Body).Decode(&rss)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return models.EmptyResponse()
	}
	return rss.ToResult()
}
