package articles

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"sort"
	"time"
)

// Handler retrieves a list of all the articles
func Handler(c *gin.Context) {
	fetcher := common.Fetcher{
		Db:         c.MustGet(common.Db).(*badger.DB),
		HttpClient: &http.Client{},
	}

	medium := Medium{Name: common.Config.MediumUsername, Fetcher: fetcher}
	wordpress := Wordpress{URL: common.Config.WPUrl, Fetcher: fetcher}
	sources := [...]common.Source{medium, wordpress}

	var allArticles = make([]Article, 0)
	for _, source := range sources {
		var articles []Article
		bytes, _ := source.GetCached()
		json.Unmarshal(bytes, &articles)
		allArticles = append(allArticles, articles...)
	}

	// Sort in descending date (i.e the most recent dates first)
	sort.Slice(allArticles, func(i, j int) bool {
		return allArticles[i].PostedAt.After(allArticles[j].PostedAt)
	})
	c.JSON(http.StatusOK, common.MakeSuccessResponse(allArticles))
}

// Article represents a single blog article
type Article struct {
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	Url       string    `json:"url"`
	PostedAt  time.Time `json:"posted_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Excerpt   string    `json:"excerpt"`
}

func getCacheKey(suffix string) string {
	return common.GetCacheKey(common.StorageArticles, suffix)
}
