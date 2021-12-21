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
	"time"
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

	var repos = make([]models.Repo, 0)
	var updatedAts = make([]time.Time, 0)

	for _, f := range fetchers {
		var result models.RepoResult
		if err := f.GetCached(&result); err != nil {
			log.Errorf("Failed to get cached data:: %v", err)
			continue
		}
		repos = append(repos, result.Repos...)
		updatedAts = append(updatedAts, result.UpdatedAt)
	}

	if len(repos) == 0 {
		c.JSON(http.StatusInternalServerError, internal.MakeErrorResponse(repos))
		return
	}

	// Sort in descending order of scores
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Score() > repos[j].Score()
	})

	if strSize != "" {
		if size, err := strconv.Atoi(strSize); err != nil {
			data := internal.MakeErrorResponse(fmt.Sprintf("%v is not a valid size", strSize))
			c.JSON(http.StatusBadRequest, data)
			return
		} else if size < len(repos) {
			index, err := getIndexOfExtra(strExtra, repos)
			extra := repos[index]
			if err == nil && index >= size { // Ensuring that the extra repo doesn't already exist in the list
				repos = repos[:(size - 1)]
				repos = append(repos, extra)
			} else {
				repos = repos[:size]
			}
		}
	}

	c.Writer.Header().Set("Cache-Control", internal.AverageCacheControl(updatedAts))
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(repos))
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
