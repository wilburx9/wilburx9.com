package common

import (
	"net/http"
	"os"
)

// HttpClient provides function fot http request
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClientMock struct {
	ResponseFilePath string
}

// Do returns an is instance http.Response with body set to the file at cm.ResponseFilePath
func (cm *HttpClientMock) Do(_ *http.Request) (*http.Response, error) {
	file, err := os.Open(cm.ResponseFilePath)
	if err != nil {
		return nil, err
	}

	return &http.Response{Body: file}, nil
}