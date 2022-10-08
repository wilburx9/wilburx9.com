package gallery

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"net/http"
)

const (
	unsplashKey   = "unsplash"
	unsplashLimit = 20
)

// Unsplash handles fetching and caching of data from Unsplash. And also returning the cached data
type Unsplash struct {
	Username   string
	AccessKey  string
	Db         database.ReadWrite
	HttpClient internal.HttpClient
}

// Cache fetches and caches Unsplash images to db
func (u Unsplash) Cache(ctx context.Context) (int, error) {
	result, err := u.fetchImages()
	if err != nil {
		return 0, err
	}

	return len(result), u.Db.Write(ctx, internal.DbGalleryKey, result...)
}

// fetchImages fetches images via HTTP
func (u Unsplash) fetchImages() ([]database.Model, error) {
	url := fmt.Sprintf("https://api.unsplash.com/users/%s/photos?per_page=%v", u.Username, unsplashLimit)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Accept-Version", "v1")
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", u.AccessKey))

	res, err := u.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil, err
	}
	defer res.Body.Close()

	var results models.UnsplashImgs
	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil, err
	}

	return results.ToImages(unsplashKey), nil
}
