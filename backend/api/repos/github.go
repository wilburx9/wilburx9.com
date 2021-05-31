package repos

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/repos/internal/models"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	githubKey = "Github"
)

// Github handles fetching and caching of Github repositories
type Github struct {
	Auth     string
	Username string
	internal.Fetch
}

// FetchAndCache fetches and saves Github repositories to DB
func (g Github) FetchAndCache() int {
	repos := g.fetchRepos()
	bytes, _ := json.Marshal(repos)
	g.CacheData(getCacheKey(githubKey), bytes)
	return len(repos)
}

// GetCached retrieves saved Github repositories
func (g Github) GetCached() ([]byte, error) {
	return g.GetCachedData(getCacheKey(githubKey))
}

func (g Github) fetchRepos() []models.Repo {
	url := "https://api.github.com/graphql"

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(getGraphQlQuery()))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return nil
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", g.Username, g.Auth)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", auth))

	resp, err := g.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil
	}
	defer resp.Body.Close()

	var data models.Github
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil
	}

	return data.ToRepos()
}

func getGraphQlQuery() string {
	queryPath := "../api/repos/internal/files/github_query.graphql"
	bytes, err := ioutil.ReadFile(queryPath)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not load graphql query file")
		return ""
	}

	// Strip it of all new line characters
	re := regexp.MustCompile(`\r?\n`)
	cleaned := re.ReplaceAllString(string(bytes), "")

	return fmt.Sprintf(`{"query":"%v"}`, cleaned)
}
