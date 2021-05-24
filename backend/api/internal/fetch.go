package internal

import (
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
)

// Fetcher has the fields needed by data fetcher structs
type Fetcher struct {
	Db         *badger.DB
	HttpClient HttpClient
}

// Source is an interface implemented all data sources. E.g gallery.instagram
type Source interface {
	FetchAndCache()
	GetCached() ([]byte, error)
}

// CacheData caches the data with key
func (f Fetcher) CacheData(key string, data []byte) {
	err := f.Db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"error": err,
		}).Error("error while saving data")
	}
}

// GetCachedData returns the data stored with key.
func (f Fetcher) GetCachedData(key string) ([]byte, error) {
	var data []byte
	err := f.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			data = val
			return nil
		})
	})
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"error": err,
		}).Error("error while getting data")
	}
	return data, err
}
