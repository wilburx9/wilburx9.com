package update

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/articles"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/repos"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var cachers []internal.Cacher

// SetUp initializes the cachers slice
func SetUp(http internal.HttpClient, db internal.Database) {
	var c = &configs.Config
	cache := internal.BaseCache{
		Db:         db,
		HttpClient: http,
	}
	instagram := gallery.Instagram{AccessToken: c.InstagramAccessToken, BaseCache: cache}
	unsplash := gallery.Unsplash{Username: c.UnsplashUsername, AccessKey: c.UnsplashAccessKey, BaseCache: cache}
	medium := articles.Medium{Name: c.MediumUsername, BaseCache: cache}
	wordpress := articles.WordPress{URL: c.WPUrl, BaseCache: cache}
	github := repos.GitHub{Auth: c.GithubToken, Username: c.UnsplashUsername, BaseCache: cache}

	cachers = []internal.Cacher{instagram, unsplash, medium, wordpress, github}
}

// Handler fetches data from sources and cache the results
func Handler(c *gin.Context, h internal.HttpClient) {
	db := c.MustGet(internal.Db).(internal.Database)
	if cap(cachers) == 0 {
		SetUp(h, db)
	}
	result := updateCache()
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(result))
}

func updateCache() map[string]interface{} {
	var startTime = time.Now()

	rc := make(chan result, len(cachers))
	var wg sync.WaitGroup

	for _, c := range cachers {
		wg.Add(1)
		go func(cacher internal.Cacher) {
			defer wg.Done()

			size, err := cacher.Cache()

			var errV *errorV
			if err != nil {
				errV = &errorV{
					Message: err.Error(),
					Details: err,
				}
			}

			rc <- result{
				Cacher: reflect.TypeOf(cacher).Name(),
				Size:   size,
				Error:  errV,
			}
		}(c)
	}
	wg.Wait()
	close(rc)

	rs := make([]result, 0, len(cachers))
	for r := range rc {
		rs = append(rs, r)
	}

	return map[string]interface{}{
		"results":  rs,
		"duration": fmt.Sprintf("%v milliseconds", time.Since(startTime).Milliseconds()),
	}
}
