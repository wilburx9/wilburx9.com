package articles

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
	"sort"
)

// Handler retrieves a list of all the articles sorted in descending creation date
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(*badger.DB),
		HttpClient: &http.Client{},
	}

	medium := Medium{Name: internal.Config.MediumUsername, Fetch: fetch}
	wordpress := Wordpress{URL: internal.Config.WPUrl, Fetch: fetch}
	fetchers := [...]internal.Fetcher{medium, wordpress}

	var allArticles = make([]models.Article, 0)
	for _, f := range fetchers {
		var articles []models.Article
		bytes, _ := f.GetCached()
		json.Unmarshal(bytes, &articles)
		allArticles = append(allArticles, articles...)
	}

	// Sort in descending date (i.e the most recent dates first)
	sort.Slice(allArticles, func(i, j int) bool {
		return allArticles[i].PostedAt.After(allArticles[j].PostedAt)
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allArticles))
}

func getCacheKey(suffix string) string {
	return internal.GetCacheKey(internal.DbArticlesKey, suffix)
}
