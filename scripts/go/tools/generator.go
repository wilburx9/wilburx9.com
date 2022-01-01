package scripts

import (
	"crypto/sha512"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
)

// GenerateHash generates a key with pass and salt
func GenerateHash(pass, salt string) string {
	hash := pbkdf2.Key([]byte(pass), []byte(salt), 200000, 50, sha512.New)
	return hex.EncodeToString(hash)
}
