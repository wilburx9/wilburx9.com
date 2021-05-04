package gallery

import (
	"encoding/json"
	"fmt"
	"github.com/wilburt/wilburx9.dev/backend"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	instagramFirestoreKey = "instagram"
	minTokenRemainingLife = 5 * time.Minute // 5 Minutes
)

type instagram struct {
	AccessToken string
	backend.Fetcher
}

func (i instagram) fetchImages() []Image {
	accessToken := i.getAccessToken()
	fields := "caption,media_type,id,media_url,timestamp,permalink,thumbnail_url,media_type"
	u := fmt.Sprintf("https://graph.instagram.com/me/media?fields=%s&access_token=%s", fields, accessToken)
	allResults := i.fetchImage([]Image{}, u)

	// Return cached result if the it fails to fetch live result
	if len(allResults) == 0 {
		return i.fetchCachedImages()
	}

	i.persistImages(allResults)
	return allResults
}
func (i instagram) persistImages(images []Image) {
	cacheKey := common.GetCacheKey(common.FirestoreGallery, instagramFirestoreKey)
	_, err := i.FsClient.Collection(common.FirestoreCache).Doc(cacheKey).Set(i.DbCtxt, images)
	if err != nil {
		common.LogError(fmt.Errorf("error while persisting instagram images to GCP Firestore: %v", err))
	}
}

func (i instagram) fetchCachedImages() []Image {
	return getImagesFrmFirestore(i.FsClient, i.DbCtxt, instagramFirestoreKey)
}

func (i instagram) fetchImage(fetched []Image, url string) []Image {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		common.LogError(err)
		return fetched
	}

	res, err := i.HttpClient.Do(req)
	if err != nil {
		common.LogError(err)
		return fetched
	}
	defer res.Body.Close()

	var data instaImgResult
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		common.LogError(err)
		return fetched
	}

	fetched = append(fetched, data.Data.toImages()...)

	// Return the fetched images if there are no more images to fetch
	if data.Paging.Next == "" {
		return fetched
	}

	return i.fetchImage(fetched, data.Paging.Next)
}

func (i instagram) getAccessToken() string {
	dSnap, err := i.FsClient.Collection(common.FirestoreTokens).Doc(instagramFirestoreKey).Get(i.DbCtxt)

	// Set the access token if it doesn't exist otherwise log an error and return an empty string
	if err != nil {
		if dSnap != nil && !dSnap.Exists() {
			return i.refreshAccessToken(i.AccessToken)
		} else {
			common.LogError(fmt.Errorf("error while fetching instagram access token %v", dSnap))
			return ""
		}
	}
	var t token
	err = dSnap.DataTo(&t)
	if err != nil {
		common.LogError(fmt.Errorf("error while unmarshalling instagram access token %v", dSnap))
		return ""
	}

	// Check if access token has expired
	if t.expired() {
		// Access token has expired. We can't refresh it
		common.LogError(fmt.Errorf("instagram access token has expired"))
		return ""
	}

	// Check if access token should be refresh
	if t.shouldRefresh() {
		return i.refreshAccessToken(t.Value)
	}
	return t.Value
}

func (i instagram) refreshAccessToken(oldToken string) string {
	u := "https://graph.instagram.com/refresh_access_token"

	params := url.Values{}
	params.Set("grant_type", "ig_refresh_token")
	params.Set("access_token", oldToken)
	payload := strings.NewReader(params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, payload)

	if err != nil {
		common.LogError(err)
		return oldToken
	}
	res, err := i.HttpClient.Do(req)
	if err != nil {
		common.LogError(err)
		return oldToken
	}
	defer res.Body.Close()

	var newT token
	err = json.NewDecoder(res.Body).Decode(&newT)
	if err != nil {
		common.LogError(err)
		return newT.Value
	}
	i.persistAccessToken(newT)
	return newT.Value
}

func (i instagram) persistAccessToken(t token) {
	_, err := i.FsClient.Collection(common.FirestoreTokens).Doc(instagramFirestoreKey).Set(i.DbCtxt, t)
	if err != nil {
		common.LogError(fmt.Errorf("error while persisting instagram access token to GCP Firestore: %v", err))
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
			Source:       "instagram",
		}
	}
	return images
}

type token struct {
	Value       string    `firestore:"value" json:"access_token"`
	ExpiresIn   int64     `firestore:"expires_in" json:"expires_in"`
	RefreshedAt time.Time `firestore:"refreshed_at,serverTimestamp"`
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
