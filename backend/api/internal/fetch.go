package internal

// Fetch has the fields needed by data fetcher structs
type Fetch struct {
	Db         Database
	HttpClient HttpClient
}

// Fetcher is an interface implemented all data fetchers
type Fetcher interface {
	FetchAndCache() int
	GetCached() ([]interface{}, error)
}