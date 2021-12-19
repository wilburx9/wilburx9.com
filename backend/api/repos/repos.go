package repos

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/repos/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
)

// Handler retrieves a list of all git repos, sorted in descending stars and forks
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(internal.Database),
		HttpClient: &http.Client{},
	}

	github := GitHub{Auth: configs.Config.GithubToken, Username: configs.Config.GithubUsername, Fetch: fetch}
	fetchers := [...]internal.Fetcher{github}

	var allRepos = make([]models.Repo, 0)
	for _, f := range fetchers {
		var result []models.Repo
		if err := f.GetCached(&result); err != nil {
			log.Errorf("Failed to get cached data:: %v", err)
			continue
		}
		allRepos = append(allRepos, result...)
	}

	if len(allRepos) == 0 {
		c.JSON(http.StatusInternalServerError, internal.MakeErrorResponse(allRepos))
		return
	}

	// Sort in descending order of scores
	sort.Slice(allRepos, func(i, j int) bool {
		return allRepos[i].Score() > allRepos[j].Score()
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allRepos))
}
