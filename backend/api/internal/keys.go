package internal

import "fmt"

const (

	// Db is the key for the Db
	Db = "storage"

	// DbGalleryKey is the key to the collection of gallery response in Db
	DbGalleryKey = "gallery"

	// DbArticlesKey is the key to the collection of articles response in Db
	DbArticlesKey = "articles"

	// DbReposKey is the key to the collection of git repos response in Db
	DbReposKey = "repositories"

	// DbKeys is the key to the collection of gateway keys
	DbKeys = "keys"
)

// GetCacheKey returns the key of cache
func GetCacheKey(group string, key string) string {
	return fmt.Sprintf("%s_%s", group, key)
}
