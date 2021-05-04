package backend

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/wilburt/wilburx9.dev/backend/common"
)

// Fetcher has the fields needed by data fetcher structs
type Fetcher struct {
	DbCtxt     context.Context
	FsClient   *firestore.Client
	HttpClient common.HttpClient
}
