package articles

import (
	"net/http"
	"os"
)

type httpClientMock struct {
	responseFilePath string
}

// Do returns an is instance http.Response with body set to the file at cm.responseFilePath
func (cm *httpClientMock) Do(_ *http.Request) (*http.Response, error) {
	file, err := os.Open(cm.responseFilePath)
	if err != nil {
		return nil, err
	}

	return &http.Response{Body: file}, nil
}