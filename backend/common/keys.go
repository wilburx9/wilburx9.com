package common

import "fmt"

const (

	// Db is the key for the Db
	Db = "storage"

	// Access is the prefix for keys that have to do with API keys, access tokens
	Access = "Access"

	// StorageCache is the key to caches of fetched data
	StorageCache = "cache"

	// StorageGallery is the key to cache of all gallery document in Db
	StorageGallery = "gallery"
)

// GetCacheKey returns the key of cache
func GetCacheKey(group string, key string) string {
	return fmt.Sprintf("%s_%s", group, key)
}
