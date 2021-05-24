package repos

import (
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"time"
)

// Handler retrieves a list of all git repos, sorted in descending stars and forks
func Handler(c *gin.Context) {

}

// repo represents a single git repository
type repo struct {
	Name        string    `json:"name"`
	Stars       int       `json:"stars"`
	Forks       int       `json:"forks"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	License     string    `json:"license"`
	Languages   []struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"languages"`
}

func getCacheKey(suffix string) string {
	return internal.GetCacheKey(internal.DbReposKey, suffix)
}
