package internal

import (
	"cloud.google.com/go/firestore"
	"context"
	log "github.com/sirupsen/logrus"
)

// Fetch has the fields needed by data fetcher structs
type Fetch struct {
	Db         *firestore.Client
	Ctx        context.Context
	HttpClient HttpClient
}

// Fetcher is an interface implemented all data fetchers
type Fetcher interface {
	FetchAndCache() int
	GetCached() ([]interface{}, error)
}

// CacheData caches the data with key
func (f Fetch) CacheData(coll string, doc string, data interface{}) {
	collection := GetDataCollection(coll)
	mapData := map[string]interface{}{"data": data}
	_, err := f.Db.Collection(collection).Doc(doc).Set(f.Ctx, mapData)
	if err != nil {
		log.Errorf("error: %v :: failed to write %q to %s.%s", err, mapData, collection, doc)
	}
}

// GetCachedData returns the data stored with key.
func (f Fetch) GetCachedData(coll string, doc string) ([]interface{}, error) {
	collection := GetDataCollection(coll)
	snapshot, err := f.Db.Collection(collection).Doc(doc).Get(f.Ctx)
	if err != nil {
		log.Errorf("Failed to read from %s.%s", collection, doc)
		return nil, err
	}

	data, err := snapshot.DataAt("data")

	if err != nil {
		log.Errorf("Failed to read snapshot data at %s.%s", collection, doc)
		return nil, err
	}
	return data.([]interface{}), nil
}
