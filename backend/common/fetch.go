package common

import (
	"github.com/dgraph-io/badger/v3"
)

// Fetcher has the fields needed by data fetcher structs
type Fetcher struct {
	Db         *badger.DB
	HttpClient HttpClient
}
