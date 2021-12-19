package repos

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/repos/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"sort"
	"strconv"
)

// Handler retrieves a list of all git repos, sorted in descending stars and forks
func Handler(c *gin.Context) {
	fetch := internal.Fetch{
		Db:         c.MustGet(internal.Db).(internal.Database),
		HttpClient: &http.Client{},
	}

	strSize := c.Query("size")   // The number of repos to return.
	strExtra := c.Query("extra") // An extra repo to add to the repos to return if strSize is less than the total repos.

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

	if strSize != "" {
		if size, err := strconv.Atoi(strSize); err != nil {
			data := internal.MakeErrorResponse(fmt.Sprintf("%v is not a valid size", strSize))
			c.JSON(http.StatusBadRequest, data)
			return
		} else if size < len(allRepos) {
			index, err := getIndexOfExtra(strExtra, allRepos)
			extra := allRepos[index]
			if err == nil && index >= size { // Ensuring that the extra repo doesn't already exist in the list
				allRepos = allRepos[:(size - 1)]
				allRepos = append(allRepos, extra)
			} else {
				allRepos = allRepos[:size]
			}
		}
	}

	c.JSON(http.StatusOK, internal.MakeSuccessResponse(allRepos))
}

func getIndexOfExtra(name string, repos []models.Repo) (int, error) {
	if name == "" {
		return 0, fmt.Errorf("name not valid")
	}
	for i := range repos {
		if repos[i].Name == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("%v not found", name)
}
