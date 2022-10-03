package update

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"github.com/wilburt/wilburx9.dev/backend/api/repos"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var cachers []database.Cacher

// SetUp initializes the cachers slice
func SetUp(http internal.HttpClient, db database.ReadWrite) {
	var c = &configs.Config

	instagram := gallery.Instagram{AccessToken: c.InstagramAccessToken, Db: db, HttpClient: http}
	unsplash := gallery.Unsplash{Username: c.UnsplashUsername, AccessKey: c.UnsplashAccessKey, Db: db, HttpClient: http}
	medium := articles.Medium{Name: c.MediumUsername, Db: db, HttpClient: http}
	wordpress := articles.WordPress{URL: c.WPUrl, Db: db, HttpClient: http}
	github := repos.GitHub{Auth: c.GithubToken, Username: c.UnsplashUsername, Db: db, HttpClient: http}

	cachers = []database.Cacher{instagram, unsplash, medium, wordpress, github}
}

// Handler fetches data from sources and cache the results
func Handler(c *gin.Context, db database.ReadWrite, h internal.HttpClient) {
	if cap(cachers) == 0 {
		SetUp(h, db)
	}
	results, duration := updateCache()

	log.Infof("Update cache: %v", generateLogMsg(results, duration))

	data := map[string]interface{}{
		"duration": duration,
	}
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(data))
}

func generateLogMsg(results []result, duration string) string {
	buffer := &bytes.Buffer{}
	for _, r := range results {
		buffer.WriteString(fmt.Sprintln("\n", "-----------------"))
		buffer.WriteString(r.Cacher)
		buffer.WriteString(fmt.Sprintln("\n", "-----------------"))

		buffer.WriteString(fmt.Sprintln("Size:", r.Size))
		if r.Error != nil {
			buffer.WriteString(fmt.Sprintln("Error:", r.Error.Message))
		}
	}
	buffer.WriteString(fmt.Sprintln("\n", "Total duration:", duration))

	return buffer.String()
}

func updateCache() ([]result, string) {
	var startTime = time.Now()

	rc := make(chan result, len(cachers))
	var wg sync.WaitGroup

	for _, c := range cachers {
		wg.Add(1)
		go func(cacher database.Cacher) {
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

	return rs, fmt.Sprintf("%v milliseconds", time.Since(startTime).Milliseconds())
}
