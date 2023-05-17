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

// AppConfig is a global container for app-wide config
var AppConfig = NewConfig()

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

// MailClient is the email marketing client
var MailClient = mailerlite.NewClient(AppConfig.MailerLiteToken)

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
