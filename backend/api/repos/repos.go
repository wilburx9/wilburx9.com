package repos

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
	"strconv"
)

const defaultLimit = 6

// Handler retrieves a list of all git repos, sorted in descending stars and forks
func Handler(c *gin.Context) {
	db := c.MustGet(internal.Db).(internal.Database)

	limit, err := getLimit(c.DefaultQuery("size", strconv.FormatInt(defaultLimit, 10)))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.MakeErrorResponse(err))
		return
	}

	repos, at, err := db.Retrieve(internal.DbReposKey, "score", limit)

	if err != nil && len(repos) == 0 {
		c.JSON(http.StatusInternalServerError, internal.MakeErrorResponse(repos))
		log.Errorf("Couldn't fetch repos. Reason :: %v", err)
		return
	}

	c.Writer.Header().Set("Cache-Control", internal.GetCacheControl(at.T))
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(repos))
}

func getLimit(strSize string) (int, error) {
	if strSize == "" {
		return defaultLimit, nil
	}

	if size, err := strconv.Atoi(strSize); err != nil || size == 0 {
		return 0, fmt.Errorf("%q is not a valid size", strSize)
	} else {
		return size, nil
	}
}
