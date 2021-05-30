package internal

import "fmt"

const (

	// Db is the key for the Db
	Db = "storage"

	// DbAccessKey is the prefix for keys that have to do with API keys, access tokens
	DbAccessKey = "DbAccessKey"

	// DbGalleryKey is the key to cache of all gallery response in Db
	DbGalleryKey = "gallery"

	// DbArticlesKey is the key to cache of all articles response in Db
	DbArticlesKey = "articles"

	// DbReposKey is the key to cache of all git repos response in Db
	DbReposKey = "repositories"
)

// GetCacheKey returns the key of cache
func GetCacheKey(group string, key string) string {
	return fmt.Sprintf("%s_%s", group, key)
}
