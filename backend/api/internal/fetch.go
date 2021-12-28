package internal

// Fetch has the fields needed by data fetcher structs
type Fetch struct {
	Db         Database
	HttpClient HttpClient
}

// Cacher is an interface that wraps functions for saving data to the db
type Cacher interface {
	Cache() int
}