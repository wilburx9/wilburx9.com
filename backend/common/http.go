package common

import "net/http"

// HttpClient provides function fot http request
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
