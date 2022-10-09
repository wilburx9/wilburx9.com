package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"net/http"
)

// HttpClient provides function fot http request
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// MakeSuccessResponse returns a template of a successful response
func MakeSuccessResponse(data interface{}) gin.H {
	return gin.H{"success": true, "data": data}
}

// MakeErrorResponse returns a template of an error response
func MakeErrorResponse(data interface{}) gin.H {
	return gin.H{"success": false, "data": data}
}

// MockHttpClient is a mock HttpClient for testing
type MockHttpClient struct {
	mock.Mock
}

// Do stub http.Do function
func (c *MockHttpClient) Do(r *http.Request) (*http.Response, error) {
	args := c.Called(r)
	response, ok := args.Get(0).(*http.Response)
	if !ok {
		response = nil
	}
	return response, args.Error(1)
}
