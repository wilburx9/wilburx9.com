package gallery

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
	"time"
)

const (
	instagramKey          = "instagram"
	tokenKey              = "token"
	minTokenRemainingLife = 24 * time.Hour * 5 // 5 Days
)

// Instagram encapsulates the fetching of Instagram images and access token management
type Instagram struct {
	AccessToken string
	internal.Fetch
}

// FetchAndCache fetches and caches data from Instagram
func (i Instagram) FetchAndCache() int {
	accessToken := i.getToken()
	fields := "caption,id,media_url,timestamp,permalink,thumbnail_url,media_type"
	u := fmt.Sprintf("https://graph.instagram.com/me/media?fields=%s&access_token=%s", fields, accessToken)
	images := i.fetchImage([]models.Image{}, u)
	i.CacheData(internal.DbGalleryKey, instagramKey, images)
	return len(images)
}

// GetCached fetches Instagram images from the db that was previously saved in Cache
func (i Instagram) GetCached() ([]interface{}, error) {
	return i.GetCachedData(internal.DbGalleryKey, instagramKey)
}

// Recursively fetch all the images
func (i Instagram) fetchImage(fetched []models.Image, url string) []models.Image {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return fetched
	}

	res, err := i.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return fetched
	}
	defer res.Body.Close()

	var data models.InstaImg
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return fetched
	}

	fetched = append(fetched, data.Data.ToImages()...)

	// Return the fetched images if there are no more images to fetch
	if data.Paging.Next == "" {
		return fetched
	}

	// Fetch the next page
	return i.fetchImage(fetched, data.Paging.Next)
}

func (i Instagram) getToken() string {
	var tk models.InstaToken
	// Attempt to get token from Db
	collection := internal.GetDataCollection(internal.DbKeys)
	snapshot, err := i.Db.Collection(collection).Doc(instagramKey).Get(i.Ctx)

	// If we haven't saved the token before, log an error and refresh the token we have now
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't fetch Instagram token")
		return i.refreshToken(i.AccessToken)
	}

	dataAt, _ := snapshot.DataAt(tokenKey)
	mapstructure.Decode(dataAt, tk)

	// Check for expired token.
	if tk.Expired() {
		// Token has expired and we can't refresh it. This should never happen.
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
	url := fmt.Sprintf("https://graph.instagram.com/refresh_access_token?grant_type=ig_refresh_token&access_token=%v", oldToken)

	req, err := http.NewRequest(http.MethodGet, url, nil)

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
	i.saveToken(newT)
	return newT.Value
}

func (i Instagram) saveToken(t models.InstaToken) {
	collection := internal.GetDataCollection(internal.DbKeys)
	_, err := i.Db.Collection(collection).Doc(instagramKey).Set(i.Ctx, t)
	if err != nil {
		log.Errorf("error while persisting Instagram access token to Db: %v", err)
	}
}
