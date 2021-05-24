package internal

import "fmt"

const (

	// Db is the key for the Db
	Db = "storage"

	// Access is the prefix for keys that have to do with API keys, access tokens
	Access = "Access"

	// StorageGallery is the key to cache of all gallery response in Db
	StorageGallery = "gallery"

	// StorageArticles is the key to cache of all articles response in Db
	StorageArticles = "articles"
)

// GetCacheKey returns the key of cache
func GetCacheKey(group string, key string) string {
	return fmt.Sprintf("%s_%s", group, key)
}
