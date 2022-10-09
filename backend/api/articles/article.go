package articles

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"net/http"
)

// Handler retrieves a list of all the articles sorted in descending creation date
func Handler(c *gin.Context, db database.ReadWrite) {
	articles, at, err := db.Read(c, internal.DbArticlesKey, "updated_on", 20)

	if err != nil || len(articles) == 0 {
		c.JSON(http.StatusInternalServerError, internal.MakeErrorResponse(articles))
		log.Errorf("Couldn't fetch articles. Reason :: %v", err)
		return
	}

	c.Writer.Header().Set("Cache-Control", internal.GetCacheControl(at.T))
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(articles))
}
