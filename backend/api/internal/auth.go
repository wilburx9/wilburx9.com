package internal

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
	"net/http"
	"strings"
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

	v := gen()
	vs := strings.Split(v, ":")

	// key := []byte(vs[0])
	auth := pbkdf2.Key([]byte(vs[0]), []byte(vs[1]), 200000, 50, sha512.New)
	return hex.EncodeToString(auth) == vs[2]
}

func gen() string {
	key := uuid.NewString()
	salt := uuid.NewString()
	hash := pbkdf2.Key([]byte(key), []byte(salt), 200000, 50, sha512.New)
	hashStr := hex.EncodeToString(hash)
	return fmt.Sprintf("%v:%v:%v", key, salt, hashStr)
}
