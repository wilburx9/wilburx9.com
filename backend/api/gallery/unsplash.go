package gallery

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
)

const (
	unsplashKey = "unsplash"
	limit       = 20
)

// Unsplash handles fetching and caching of data from Unsplash. And also returning the cached data
type Unsplash struct {
	Username  string
	AccessKey string
	internal.Fetch
}

// Cache fetches and caches Unsplash images to db
func (u Unsplash) Cache() int {
	result := u.FetchImages()
	err := u.Db.Persist(internal.DbGalleryKey, result)
	if err != nil {
		log.Errorf("Couldn't cache Unsplash images. Reason :: %v", err)
		return 0
	}
	return len(result)
}

// FetchImages fetches images via HTTP
func (u Unsplash) FetchImages() []internal.DbModel {
	url := fmt.Sprintf("https://api.unsplash.com/users/%s/photos?per_page=%v", u.Username, limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return nil
	}

	req.Header.Add("Accept-Version", "v1")
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", u.AccessKey))

	res, err := u.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil
	}
	defer res.Body.Close()

	var results models.UnsplashImgs
	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil
	}

	return results.ToImages(unsplashKey)
}
