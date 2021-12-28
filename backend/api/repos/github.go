package repos

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/repos/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	githubKey = "github"
)

// GitHub handles fetching and caching of GitHub repositories
type GitHub struct {
	Auth     string
	Username string
	internal.BaseCache
}

// Cache fetches and saves GitHub repositories to DB
func (g GitHub) Cache() (int, error) {
	result, err := g.fetchRepos()
	if err != nil {
		return 0, err
	}

	return len(result), g.Db.Persist(internal.DbReposKey, result...)
}

func (g GitHub) fetchRepos() ([]internal.DbModel, error) {
	url := "https://api.github.com/graphql"

	query, err := getGraphQlQuery()
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(query))
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", g.Username, g.Auth)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", auth))

	resp, err := g.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil, err
	}
	defer resp.Body.Close()

	var data models.GitHub
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil, err
	}

	return data.ToResult(githubKey), nil
}

func getGraphQlQuery() (string, error) {
	queryPath := fmt.Sprintf("%v/api/repos/internal/files/github_query.graphql", configs.Config.AppHome)
	bytes, err := ioutil.ReadFile(queryPath)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not load graphql query file")
		return "", err
	}

	// Strip it of all new line characters
	re := regexp.MustCompile(`\r?\n`)
	cleaned := re.ReplaceAllString(string(bytes), "")

	return fmt.Sprintf(`{"query":"%v"}`, cleaned), nil
}
