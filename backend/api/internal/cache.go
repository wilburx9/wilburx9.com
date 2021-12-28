package internal

// BaseCache has the fields needed by data fetcher structs
type BaseCache struct {
	Db         Database
	HttpClient HttpClient
}

// Cacher is an interface that wraps functions for saving data to the db
type Cacher interface {
	Cache() (int, error)
}
