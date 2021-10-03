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
	articles := m.fetchArticles()
	m.CacheData(internal.DbArticlesKey, mediumKey, articles)
	return len(articles)
}

// GetCached returns cached Medium articles
func (m Medium) GetCached() ([]interface{}, error) {
	return m.GetCachedData(internal.DbArticlesKey, mediumKey)
}

// fetchArticles fetches articles via HTTP
func (m Medium) fetchArticles() []models.Article {
	url := fmt.Sprintf("https://medium.com/feed/%s", m.Name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return nil
	}

	res, err := m.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil
	}
	defer res.Body.Close()

	var rss models.Rss
	err = xml.NewDecoder(res.Body).Decode(&rss)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil
	}
	return rss.ToArticles()
}
