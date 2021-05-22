package gallery

import (
	"encoding/json"
	"fmt"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"strconv"
)

const (
	unsplashKey = "Unsplash"
)

// Unsplash handles fetching and caching of data from Unsplash. And also returning the cached data
type Unsplash struct {
	Username  string
	AccessKey string
	common.Fetcher
}

// FetchAndCache fetches data from Unsplash and caches it
func (u Unsplash) FetchAndCache() {
	images := u.fetchImage([]Image{}, 1)
	buf, _ := json.Marshal(images)
	u.CacheData(getCacheKey(unsplashKey), buf)
}

// GetCached returns data that was cached in Cache
func (u Unsplash) GetCached() ([]byte, error) {
	return u.GetCachedData(getCacheKey(unsplashKey))
}

func (u Unsplash) fetchImage(fetched []Image, page int) []Image {
	url := fmt.Sprintf("https://api.Unsplash.com/users/%s/photos?page=%d&per_page=5", u.Username, page) // TODO: Increment per_page to 30 after testing this
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		common.LogError(err)
		return fetched
	}

	req.Header.Add("Accept-Version", "v1")
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", u.AccessKey))

	res, err := u.HttpClient.Do(req)
	if err != nil {
		common.LogError(err)
		return fetched
	}
	defer res.Body.Close()

	var results unsplashImgSlice
	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		common.LogError(err)
		return fetched
	}

	fetched = append(fetched, results.toImages()...)

	totalImages, err := strconv.Atoi(res.Header.Get("X-Total"))

	// Return if an error is encountered or if all images has been fetched
	if err != nil || len(fetched) >= totalImages {
		return fetched
	}

	page++
	return u.fetchImage(fetched, page)
}

func (m unsplashImgSlice) toImages() []Image {
	var timeLayout = "2006-01-02T03:04:05-07:00"
	var images = make([]Image, len(m))

	for i, e := range m {
		images[i] = Image{
			SrcThumbnail: e.Urls.Small,
			Src:          e.Urls.Full,
			Url:          e.Links.HTML,
			Caption:      e.Description,
			UploadedAt:   common.StringToTime(timeLayout, e.CreatedAt),
			Source:       "Unsplash",
			Meta: map[string]interface{}{
				"user": e.User,
			},
		}
	}
	return images
}

type unsplashImgSlice []unsplashImg

type unsplashImg struct {
	CreatedAt   string `json:"created_at"`
	Color       string `json:"color"`
	Description string `json:"description"`
	User        user   `json:"user"`
	Urls        struct {
		Full  string `json:"full"`
		Small string `json:"small"`
	} `json:"urls"`
	Links struct {
		HTML string `json:"html"`
	} `json:"links"`
}

type user struct {
	Username string `json:"Username"`
	Name     string `json:"name"`
}
