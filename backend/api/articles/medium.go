package articles

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"net/http"
)

const (
	mediumKey = "medium"
)

// Medium encapsulates the fetching and caching of medium articles
type Medium struct {
	Name string // should be Medium username (e.g "@Wilburx9") or publication (e.g. flutter-community)
	Db         database.ReadWrite
	HttpClient internal.HttpClient
}

// Cache fetches and caches all Medium Articles
func (m Medium) Cache() (int, error) {
	result, err := m.fetchArticles()
	if err != nil {
		return 0, err
	}

	return len(result), m.Db.Write(internal.DbArticlesKey, result...)
}

// fetchArticles fetches articles via HTTP
func (m Medium) fetchArticles() ([]database.Model, error) {
	url := fmt.Sprintf("https://medium.com/feed/%s", m.Name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	res, err := m.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil, err
	}
	defer res.Body.Close()

	var rss models.Rss
	err = xml.NewDecoder(res.Body).Decode(&rss)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil, err
	}
	return rss.ToResult(mediumKey), nil
}
