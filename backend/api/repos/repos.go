package repos

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/repos/internal/models"
	"net/http"
	"sort"
)

// Handler retrieves a list of all git repos, sorted in descending stars and forks
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(*badger.DB),
		HttpClient: &http.Client{},
	}

	github := Github{Auth: internal.Config.GithubToken, Username: internal.Config.GithubUsername, Fetch: fetch}
	fetchers := [...]internal.Fetcher{github}

	var allRepos = make([]models.Repo, 0)
	for _, f := range fetchers {
		var repos []models.Repo
		bytes, _ := f.GetCached()
		json.Unmarshal(bytes, &repos)
		allRepos = append(allRepos, repos...)
	}

	// Sort in descending order of scores
	sort.Slice(allRepos, func(i, j int) bool {
		return allRepos[i].Score() > allRepos[j].Score()
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allRepos))
}

func getCacheKey(suffix string) string {
	return internal.GetCacheKey(internal.DbReposKey, suffix)
}
