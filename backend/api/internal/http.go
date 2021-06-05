package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	Header           http.Header
}

// Do returns an is instance http.Response with body set to the file at cm.ResponseFilePath
func (cm *HttpClientMock) Do(_ *http.Request) (*http.Response, error) {
	file, err := os.Open(cm.ResponseFilePath)
	if err != nil {
		log.Errorln(fmt.Sprintf("error while opening file %v", err))
		return nil, err
	}

	return &http.Response{Body: file, Header: cm.Header}, nil
}

// MakeSuccessResponse returns a template of a successful response
func MakeSuccessResponse(data interface{}) gin.H {
	return gin.H{"success": true, "data": data}
}

// MakeErrorResponse returns a template of an error response
func MakeErrorResponse(data interface{}) gin.H {
	return gin.H{"success": false, "data": data}
}
