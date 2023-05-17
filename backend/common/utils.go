package common

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mailerlite/mailerlite-go"
	"net/http"
	"strings"
	"time"
)

const (
	Photography = "photography"
	Programming = "programming"
	Blog        = "blog"
)

var AppConfig = NewConfig()

func getResponseBody(success bool, data any) string {

	res := map[string]any{
		"success": success,
		"data":    data,
	}

	b, err := json.Marshal(res)
	if err != nil {
		return "Something went wrong"
	}
	return string(b)
}

// MakeResponse returns an error or success lambda response
func MakeResponse(origin string, code int, data any) events.APIGatewayProxyResponse {
	success := false
	if code >= 200 && code <= 299 {
		success = true
	}

	headers := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Methods": "OPTIONS,POST",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
	}
	for _, o := range AppConfig.AllowedOrigins {
		if strings.EqualFold(o, origin) {
			headers["Access-Control-Allow-Origin"] = origin
			break
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       getResponseBody(success, data),
		Headers:    headers,
	}
}

// HttpClient returns an instance of Http Client with custom config
var HttpClient = &http.Client{
	Timeout: time.Second * 20,
}

var MailClient = mailerlite.NewClient(AppConfig.MailerLiteToken)
