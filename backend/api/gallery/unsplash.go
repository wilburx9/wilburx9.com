package gallery

// import (
// 	"encoding/json"
// 	"fmt"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
// 	"github.com/wilburt/wilburx9.dev/backend/api/internal"
// 	"net/http"
// 	"strconv"
// )
//
// const (
// 	unsplashKey = "unsplash"
// )
//
// // Unsplash handles fetching and caching of data from Unsplash. And also returning the cached data
// type Unsplash struct {
// 	Username  string
// 	AccessKey string
// 	internal.Fetch
// }
//
// // FetchAndCache fetches data from Unsplash and caches it
// func (u Unsplash) FetchAndCache() int {
// 	images := u.FetchImage([]models.Image{}, 1)
// 	u.CacheData(internal.DbGalleryKey, unsplashKey, images)
// 	return len(images)
// }
//
// // GetCached returns data that was cached in Cache
// func (u Unsplash) GetCached() ([]interface{}, error) {
// 	return u.GetCachedData(internal.DbGalleryKey, unsplashKey)
// }
//
// // FetchImage fetches images via HTTP
// func (u Unsplash) FetchImage(fetched []models.Image, page int) []models.Image {
// 	url := fmt.Sprintf("https://api.Unsplash.com/users/%s/photos?page=%d&per_page=5", u.Username, page) // TODO: Increment per_page to 30 after testing this
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
// 		return fetched
// 	}
//
// 	req.Header.Add("Accept-Version", "v1")
// 	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", u.AccessKey))
//
// 	res, err := u.HttpClient.Do(req)
// 	if err != nil {
// 		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
// 		return fetched
// 	}
// 	defer res.Body.Close()
//
// 	var results models.UnsplashImgs
// 	err = json.NewDecoder(res.Body).Decode(&results)
// 	if err != nil {
// 		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
// 		return fetched
// 	}
//
// 	fetched = append(fetched, results.ToImages()...)
//
// 	totalImages, err := strconv.Atoi(res.Header.Get("X-Total"))
//
// 	// Return if an error is encountered or if all images has been fetched
// 	if err != nil || len(fetched) >= totalImages {
// 		return fetched
// 	}
//
// 	page++
// 	return u.FetchImage(fetched, page)
// }
