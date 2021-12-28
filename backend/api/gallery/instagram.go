package gallery

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"math"
	"net/http"
	"net/url"
	"time"
)

const (
	instagramKey          = "instagram"
	minTokenRemainingLife = 24 * time.Hour * 5 // 5 Days
	instagramLimit        = "20"
)

// Instagram encapsulates the fetching of Instagram images and access token management
type Instagram struct {
	AccessToken string
	internal.Fetch
}

// Cache fetches and caches Instagram images to db
func (i Instagram) Cache() int {
	result := i.fetchImages()
	err := i.Db.Persist(internal.DbGalleryKey, result...)
	if err != nil {
		log.Errorf("Couldn't cache Instagram images. Reason :: %v", err)
		return 0
	}
	return len(result)
}

// Recursively fetch all the images
func (i Instagram) fetchImages() []internal.DbModel {
	u, _ := url.Parse("https://graph.instagram.com/me/media")
	u.Query().Set("fields", "caption,id,media_url,timestamp,permalink,thumbnail_url,media_type")
	u.Query().Set("access_token", i.getToken())
	u.Query().Set("limit", instagramLimit)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return nil
	}

	res, err := i.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil
	}
	defer res.Body.Close()

	var data models.InstaImg
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil
	}

	return data.Data.ToImages(instagramKey)
}

func (i Instagram) getToken() string {

	// Attempt to get token from Db
	keys, _, err := i.Db.Retrieve(internal.DbKeys, "", math.MaxInt)
	// If we haven't saved the token before, log an error and refresh the token we have now
	if err != nil {
		return i.refreshToken(i.AccessToken)
	}

	var tk models.InstaToken
	// Get token map from keys
	for _, m := range keys {
		if val, ok := m[instagramKey]; ok {
			if bytes, err := json.Marshal(val); err != nil {
				json.Unmarshal(bytes, &tk)
			}
			break
		}
	}

	// Check for expired token.
	if tk.Expired() {
		// Token has expired and it can't be refreshed. This should never happen.
		log.Error("Instagram access token has expired")
		return ""
	}

	// Refresh the token if need be
	if tk.ShouldRefresh(minTokenRemainingLife) {
		return i.refreshToken(tk.Value)
	}
	return tk.Value
}

func (i Instagram) refreshToken(oldToken string) string {

	u, _ := url.Parse("https://graph.instagram.com/refresh_access_token")
	u.Query().Set("grant_type", "ig_refresh_token")
	u.Query().Set("access_token", oldToken)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return oldToken
	}
	res, err := i.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return oldToken
	}
	defer res.Body.Close()

	var newT models.InstaToken
	err = json.NewDecoder(res.Body).Decode(&newT)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall refresh token response")
		return newT.Value
	}
	newT.RefreshedAt = time.Now()
	err = i.Db.Persist(internal.DbKeys, newT)
	if err != nil {
		log.Errorf("Couldn't save Instagram token. Reason :: %v", err)
	}
	return newT.Value
}
