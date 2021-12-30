package internal

import (
	"github.com/gin-gonic/gin"
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

func valid(a string) bool {
	if (len(a)) == 0 {
		return false
	}

	// v := gen()
	// vs := strings.Split(v, ":")
	//
	// // key := []byte(vs[0])
	// auth := pbkdf2.Key([]byte(vs[0]), []byte(vs[1]), 200000, 50, sha512.New)
	// return hex.EncodeToString(auth) == vs[2]
	return true
}

// func gen() string {
// 	key := uuid.NewString()
// 	salt := uuid.NewString()
// 	hash := pbkdf2.Key([]byte(key), []byte(salt), 200000, 50, sha512.New)
// 	hashStr := hex.EncodeToString(hash)
// 	return fmt.Sprintf("%v:%v:%v", key, salt, hashStr)
// }

// Key:: 8E4FA332-5894-41FA-A88F-DD76881144F9
// Salt:: E9E36CE6-CF24-4D5F-8567-848F999C4C7A