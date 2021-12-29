package gallery

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"net/http"
)

// Handler retrieves a list of all the images sorted in descending creation date
func Handler(c *gin.Context) {
	db := c.MustGet(internal.Db).(database.ReadWrite)
	images, at, err := db.Read(internal.DbGalleryKey, "uploaded_on", 30)


	if err != nil && len(images) == 0 {
		c.JSON(http.StatusInternalServerError, internal.MakeErrorResponse(images))
		log.Errorf("Couldn't fetch gallery. Reason :: %v", err)
		return
	}

	c.Writer.Header().Set("Cache-Control", internal.GetCacheControl(at.T))
	c.JSON(http.StatusOK, internal.MakeSuccessResponse(images))
}