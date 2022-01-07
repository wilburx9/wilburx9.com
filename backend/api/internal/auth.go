package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	scripts "github.com/wilburt/wilburx9.dev/scripts/tools"
	"net/http"
)

// AuthMiddleware validates that the request has a valid authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if valid(c.GetHeader("Authorization")) {
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, MakeErrorResponse("Authorization missing or invalid"))
		c.Abort()
	}
}

func valid(pass string) bool {
	if configs.Config.IsDebug() {
		return true
	}
	// Since we are using UUIDs, the expected length is 36 runes.
	// Bounding it with >= 1 and <=49 just in case.
	if len(pass) > 0 && len(pass) < 50 {
		return configs.Config.APIHash == scripts.GenerateHash(pass, configs.Config.APISalt)
	}
	return false
}
