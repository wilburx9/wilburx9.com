package update

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/articles"
	"github.com/wilburt/wilburx9.dev/backend/api/email"
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
func Handler(c *gin.Context, h internal.HttpClient) {
	db := c.MustGet(internal.Db).(database.ReadWrite)
	if cap(cachers) == 0 {
		SetUp(h, db)
	}
	results, duration := updateCache()

	err := email.Send(generateEmail(results, duration), h)

	data := map[string]interface{}{
		"results":              results,
		"duration":             duration,
		"email_report_success": err == nil,
	}
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(data))
}

func generateEmail(results []result, duration string) email.Data {
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

	return email.Data{
		SenderEmail: configs.Config.EmailReceiver,
		SenderName:  "Jesse Bruce Pinkman",
		Subject:     "Yo! Batch cache update report!",
		Message:     buffer.String(),
	}
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
