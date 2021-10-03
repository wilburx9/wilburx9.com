package repos

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/repos/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
)

// Handler retrieves a list of all git repos, sorted in descending stars and forks
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(*firestore.Client),
		Ctx:        c,
		HttpClient: &http.Client{},
	}

	github := GitHub{Auth: configs.Config.GithubToken, Username: configs.Config.GithubUsername, Fetch: fetch}
	fetchers := [...]internal.Fetcher{github}

	var allRepos = make([]models.Repo, 0)
	for _, f := range fetchers {
		maps, _ := f.GetCached()
		allRepos = append(allRepos, mapSliceTRepos(maps)...)
	}

	// Sort in descending order of scores
	sort.Slice(allRepos, func(i, j int) bool {
		return allRepos[i].Score() > allRepos[j].Score()
	})
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allRepos))
}

func mapSliceTRepos(maps []interface{}) []models.Repo {
	var repos = make([]models.Repo, len(maps))
	for i, v := range maps {
		var repo models.Repo
		err := mapstructure.Decode(v, &repo)
		if err == nil {
			repos[i] = repo
		}
	}
	return repos
}
