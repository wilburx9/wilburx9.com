package common

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mailerlite/mailerlite-go"
	"log"
	"net/http"
	"strings"
	"time"
)

var Groups = []string{"photography", "software"}

const Blog = "blog"

func init() {
	config, err := newConfig()
	if err != nil {
		log.Println(err)
		return
	}

	AppConfig = config
	MailClient = mailerlite.NewClient(AppConfig.MailerLiteToken)
}

// AppConfig is a global container for app-wide config
var AppConfig *Config

// MailClient is the email marketing client
var MailClient *mailerlite.Client

// InitSuccess returns true if global variables have successful initialized
func InitSuccess() bool {
	return AppConfig != nil && MailClient != nil
}

// HttpClient returns an instance of Http Client with custom config
var HttpClient = &http.Client{Timeout: time.Second * 20}

// GenerateResponse returns an error or success lambda response
func GenerateResponse(origin string, code int, data any) events.APIGatewayProxyResponse {
	success := false
	if code >= 200 && code <= 299 {
		success = true
	}

	headers := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Methods": "OPTIONS,POST",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
	}
	if AppConfig != nil {
		for _, o := range AppConfig.AllowedOrigins {
			if strings.EqualFold(o, origin) {
				headers["Access-Control-Allow-Origin"] = origin
				break
			}
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       getResponseBody(success, data),
		Headers:    headers,
	}
}

func getResponseBody(success bool, data any) string {

	res := map[string]any{
		"success": success,
		"data":    data,
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Println("error while marshalling response body: ", err)
		return "Something went wrong"
	}
	return string(b)
}
