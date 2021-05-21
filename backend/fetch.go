package backend

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/wilburt/wilburx9.dev/backend/common"
)

// Fetcher has the fields needed by data fetcher structs
type Fetcher struct {
	Db         *badger.DB
	HttpClient common.HttpClient
}
