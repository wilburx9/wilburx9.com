package common

import (
	"net/http"
	"os"
)

// HttpClient provides function fot http request
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HttpClientMock is a mock HttpClient for testing
type HttpClientMock struct {
	ResponseFilePath string
	Header http.Header
}

// Do returns an is instance http.Response with body set to the file at cm.ResponseFilePath
func (cm *HttpClientMock) Do(_ *http.Request) (*http.Response, error) {
	file, err := os.Open(cm.ResponseFilePath)
	if err != nil {
		return nil, err
	}

	return &http.Response{Body: file, Header: cm.Header}, nil
}
