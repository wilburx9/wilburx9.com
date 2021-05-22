package gallery

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	instagramKey          = "Instagram"
	minTokenRemainingLife = 5 * time.Minute // 5 Minutes
)

// Instagram encapsulates the fetching data from Instagram, caching the data,
// fetching cached data, and refreshing Instagram access token
type Instagram struct {
	AccessToken string
	common.Fetcher
}

// FetchAndCache fetches and caches data fetched from Instagram
func (i Instagram) FetchAndCache() {
	accessToken := i.getToken()
	fields := "caption,media_type,id,media_url,timestamp,permalink,thumbnail_url,media_type"
	u := fmt.Sprintf("https://graph.Instagram.com/me/media?fields=%s&access_token=%s", fields, accessToken)
	allResults := i.fetchImage([]Image{}, u)
	buf, _ := json.Marshal(allResults)
	i.CacheData(getCacheKey(instagramKey), buf)
}

// GetCached fetches Instagram images from the db that was previously saved in Cache
func (i Instagram) GetCached() ([]byte, error) {
	return i.GetCachedData(getCacheKey(instagramKey))
}

// Recursively fetch all the images
func (i Instagram) fetchImage(fetched []Image, url string) []Image {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Warning(err)
		return fetched
	}

	res, err := i.HttpClient.Do(req)
	if err != nil {
		log.Warning(err)
		return fetched
	}
	defer res.Body.Close()

	var data instaImgResult
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Error(err)
		return fetched
	}

	fetched = append(fetched, data.Data.toImages()...)

	// Return the fetched images if there are no more images to fetch
	if data.Paging.Next == "" {
		return fetched
	}

	// Fetch the next page
	return i.fetchImage(fetched, data.Paging.Next)
}

func (i Instagram) getToken() string {
	var tk token
	// Attempt to get token from Db
	err := i.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(getInstagramToken()))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &tk)
		})
	})

	// If we haven't saved the token before, log and error and refresh the token we have
	if err != nil {
		log.Warningf("error while fetching Instagram access token %v", err)
		return i.refreshToken(i.AccessToken)
	}

	// Check for expired token. This shouldn't happen normally
	if tk.expired() {
		// Access token has expired. We can't refresh it
		log.Error("instagram access token has expired")
		return ""
	}

	// Refresh the token if needs be
	if tk.shouldRefresh() {
		return i.refreshToken(tk.Value)
	}
	return tk.Value
}

func (i Instagram) refreshToken(oldToken string) string {
	u := "https://graph.Instagram.com/refresh_access_token"

	params := url.Values{}
	params.Set("grant_type", "ig_refresh_token")
	params.Set("access_token", oldToken)
	payload := strings.NewReader(params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, payload)

	if err != nil {
		log.Warning(err)
		return oldToken
	}
	res, err := i.HttpClient.Do(req)
	if err != nil {
		log.Warning(err)
		return oldToken
	}
	defer res.Body.Close()

	var newT token
	err = json.NewDecoder(res.Body).Decode(&newT)
	if err != nil {
		log.Error(err)
		return newT.Value
	}
	newT.RefreshedAt = time.Now()
	i.saveToken(newT)
	return newT.Value
}

func (i Instagram) saveToken(t token) {
	err := i.Db.Update(func(txn *badger.Txn) error {
		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return txn.Set([]byte(getInstagramToken()), buf)
	})
	if err != nil {
		log.Errorf("error while persisting Instagram access token to Db: %v", err)
	}
}

func (t token) expired() bool {
	var now = time.Now()
	var expireTime = t.RefreshedAt.Add(time.Minute * time.Duration(t.ExpiresIn))
	return now.Equal(expireTime) || now.After(expireTime)
}

// It should be refreshed if the remaining life of the access token is less than 5 days
func (t token) shouldRefresh() bool {
	var now = time.Now()
	var expireTime = t.RefreshedAt.Add(time.Minute * time.Duration(t.ExpiresIn))
	diff := expireTime.Sub(now)
	return diff <= minTokenRemainingLife
}

func getInstagramToken() string {
	return common.GetCacheKey(common.Access, instagramKey)
}

func (s instaImgSlice) toImages() []Image {
	var timeLayout = "2006-02-01T15:04:05-0700"
	var images = make([]Image, len(s))

	for i, e := range s {
		images[i] = Image{
			SrcThumbnail: e.MediaURL,
			Src:          e.MediaURL,
			Url:          e.Permalink,
			Caption:      e.Caption,
			UploadedAt:   common.StringToTime(timeLayout, e.Timestamp),
			Source:       "Instagram",
		}
	}
	return images
}

type token struct {
	Value       string    `json:"access_token"`
	ExpiresIn   int64     `json:"expires_in"`
	RefreshedAt time.Time `json:"refreshed_at"`
}

type instaImgSlice []instaImg

type instaImgResult struct {
	Data   instaImgSlice `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

type instaImg struct {
	Caption   string `json:"caption"`
	MediaType string `json:"media_type"`
	ID        string `json:"id"`
	MediaURL  string `json:"media_url"`
	Timestamp string `json:"timestamp"`
	Permalink string `json:"permalink"`
}
