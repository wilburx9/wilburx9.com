package gallery

import (
	"encoding/json"
	"fmt"
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
	token, err := i.getToken()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't get token")
		return nil
	}

	u, _ := url.Parse("https://graph.instagram.com/me/media")
	q := u.Query()
	q.Set("fields", "caption,id,media_url,timestamp,permalink,thumbnail_url,media_type")
	q.Set("access_token", token)
	q.Set("limit", instagramLimit)
	u.RawQuery = q.Encode()

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

func (i Instagram) getToken() (string, error) {

	// Attempt to get token from Db
	keys, _, err := i.Db.Retrieve(internal.DbKeys, "", math.MaxInt)
	// If we haven't saved the token before, log an error and refresh the token we have now
	if err != nil || len(keys) == 0 {
		log.Warningf("Couldn't get Keys from the db. Possible error: %v", err)
		return i.refreshToken(i.AccessToken)
	}

	var tk models.InstaToken
	// Get token map from keys
	for _, m := range keys {
		if _, ok := m["id"]; ok && m["id"] == instagramKey {
			if bytes, err := json.Marshal(m); err == nil {
				json.Unmarshal(bytes, &tk)
			}
			break
		}
	}

	if tk.ID == "" || tk.Value == "" {
		return i.refreshToken(i.AccessToken)
	}

	// Check for expired token.
	if tk.Expired() {
		// Token has expired and it can't be refreshed. This should never happen.
		return "", fmt.Errorf("instagram access token has expired")
	}

	// Refresh the token if need be
	if tk.ShouldRefresh(minTokenRemainingLife) {
		return i.refreshToken(tk.Value)
	}
	return tk.Value, nil
}

func (i Instagram) refreshToken(oldToken string) (string, error) {
	u, _ := url.Parse("https://graph.instagram.com/refresh_access_token")
	q := u.Query()
	q.Set("grant_type", "ig_refresh_token")
	q.Set("access_token", oldToken)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)

	res, err := i.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var newT = models.NewInstaToken(instagramKey)
	err = json.NewDecoder(res.Body).Decode(&newT)
	if err != nil {
		return "", err
	}

	newT.RefreshedAt = time.Now()
	err = i.Db.Persist(internal.DbKeys, newT)
	if err != nil {
		return "", err
	}

	return newT.Value, nil
}
