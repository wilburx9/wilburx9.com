package database

import (
	"context"
	"time"
)

// Model wraps the id method
type Model interface {
	Id() string
}

// UpdatedAt is a container that holds the last time a collection was updated
type UpdatedAt struct {
	T time.Time `firestore:"updated_at,serverTimestamp"`
}

// ReadWrite is the interface that wraps basic functions for reading and writing to and from concrete databases
type ReadWrite interface {
	Read(ctx context.Context, source, orderBy string, limit int) ([]map[string]interface{}, UpdatedAt, error)
	Write(ctx context.Context, source string, models ...Model) error
	Close()
}

// Cacher is an interface that wraps the function for fetching data from APIs (e.g. Instagram, GitHub etc.)
type Cacher interface {
	Cache(ctx context.Context) (int, error)
}
