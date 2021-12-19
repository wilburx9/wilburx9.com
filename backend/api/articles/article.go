package articles

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
)

// Handler retrieves a list of all the articles sorted in descending creation date
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(internal.Database),
		HttpClient: &http.Client{},
	}

	medium := Medium{Name: configs.Config.MediumUsername, Fetch: fetch}
	wordpress := WordPress{URL: configs.Config.WPUrl, Fetch: fetch}
	fetchers := [...]internal.Fetcher{medium, wordpress}

	var allArticles = make([]models.Article, 0)
	for _, f := range fetchers {
		var result []models.Article
		if err := f.GetCached(&result); err != nil {
			log.Errorf("Failed to get cached data:: %v", err)
			continue
		}
		allArticles = append(allArticles, result...)
	}

	if len(allArticles) == 0 {
		c.JSON(http.StatusInternalServerError, internal.MakeErrorResponse(allArticles))
		return
	}

	// Sort in descending date (i.e. the most recent dates first)
	sort.Slice(allArticles, func(i, j int) bool {
		return allArticles[i].PostedAt.After(allArticles[j].PostedAt)
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allArticles))
}
