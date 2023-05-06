package common

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"time"
)

func getResponseBody(success bool, data interface{}) string {

	res := map[string]interface{}{
		"success": success,
		"data":    data,
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return "Something went wrong"
	}
	return string(bytes)
}

// MakeResponse returns an error or success lambda response
func MakeResponse(code int, data interface{}) events.APIGatewayProxyResponse {
	success := false
	if code <= 299 {
		success = true
	}
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       getResponseBody(success, data),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// HttpClient returns an instance of Http Client with custom config
var HttpClient = &http.Client{
	Timeout: time.Second * 20,
}
