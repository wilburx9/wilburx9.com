package common

import "fmt"

const (

	// StorageCtxt is the key for storage context
	StorageCtxt = "storage/context"

	// StorageFirestore is the key fore Firestore
	StorageFirestore = "storage/GCPFirestore"

	// FirestoreTokens is the key to the tokens collection on GCP Firestore
	FirestoreTokens = "tokens"

	// FirestoreCache is the key to caches of fetched data
	FirestoreCache = "cache"

	// FirestoreGallery is the key to cache of all gallery document in GCP Firestore
	FirestoreGallery = "gallery"
)

// GetCacheKey returns the key of cache
func GetCacheKey(group string, key string) string {
	return fmt.Sprintf("%s_%s", group, key)
}
