package articles

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
)

// Handler retrieves a list of all the articles sorted in descending creation date
func Handler(c *gin.Context) {
	// TODO: Use query params to for list size, and add special consideration for a repo
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(*firestore.Client),
		Ctx:        c,
		HttpClient: &http.Client{},
	}

	medium := Medium{Name: configs.Config.MediumUsername, Fetch: fetch}
	wordpress := WordPress{URL: configs.Config.WPUrl, Fetch: fetch}
	fetchers := [...]internal.Fetcher{medium, wordpress}

	var allArticles = make([]models.Article, 0)
	for _, f := range fetchers {
		maps, _ := f.GetCached()
		allArticles = append(allArticles, mapSliceToArticles(maps)...)
	}

	// Sort in descending date (i.e. the most recent dates first)
	sort.Slice(allArticles, func(i, j int) bool {
		return allArticles[i].PostedAt.After(allArticles[j].PostedAt)
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allArticles))
}

func mapSliceToArticles(maps []interface{}) []models.Article {
	var articles = make([]models.Article, len(maps))
	for i, v := range maps {
		var article models.Article
		err := mapstructure.Decode(v, &article)
		if err == nil {
			articles[i] = article
		}
	}
	return articles
}
