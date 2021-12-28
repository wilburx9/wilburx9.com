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
	"os"
	"path/filepath"
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
	internal.Fetch
}

// Cache fetches and saves GitHub repositories to DB
func (g GitHub) Cache() int {
	result := g.fetchRepos()
	err := g.Db.Persist(internal.DbReposKey, result...)
	if err != nil {
		log.Errorf("Couldn't cache Github repos. Reason :: %v", err)
		return 0
	}
	return len(result)
}

func (g GitHub) fetchRepos() []internal.DbModel {
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

	var data models.GitHub
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return nil
	}

	return data.ToResult(githubKey)
}

func getGraphQlQuery() string {
	queryPath := fmt.Sprintf("%v/api/repos/internal/files/github_query.graphql", configs.Config.AppHome)
	bytes, err1 := ioutil.ReadFile(queryPath)
	if err1 != nil {
		root, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatalf("Could not read root directory %s", err)
			return ""
		}
		log.Println(fmt.Sprintf("Current directory is %v", root))
		log.WithFields(log.Fields{"error": err1}).Error("Could not load graphql query file")
		return ""
	}

	// Strip it of all new line characters
	re := regexp.MustCompile(`\r?\n`)
	cleaned := re.ReplaceAllString(string(bytes), "")

	return fmt.Sprintf(`{"query":"%v"}`, cleaned)
}
