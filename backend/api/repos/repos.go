package repos

import (
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
)

// Handler retrieves a list of all git repos, sorted in descending stars and forks
func Handler(c *gin.Context) {

}

func getCacheKey(suffix string) string {
	return internal.GetCacheKey(internal.DbReposKey, suffix)
}
