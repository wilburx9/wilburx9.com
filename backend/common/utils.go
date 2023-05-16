package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	Photography = "photography"
	Programming = "programming"
	Blog        = "blog"
)

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
func MakeResponse(code int, data any) events.APIGatewayProxyResponse {
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

// MakeMailChimpRequest makes http calls to MailChimp API
func MakeMailChimpRequest(ctx context.Context, method string, path string, reqBody any, respBody any) bool {
	config := ConfigFromContext(&ctx)
	u := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/%s", config.MailChimpDC, path)

	reqBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBuffer).Encode(reqBody)
	if err != nil {
		log.Println("error while encoding request body: ", err)
		return false
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBuffer)
	if err != nil {
		log.Println("error while creating http request: ", err)
		return false
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.MailChimpToken))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("request returned an error: ", err)
		return false
	}

	defer res.Body.Close()

	// Any non 2xx response code denotes an error. See https://mailchimp.com/developer/marketing/docs/errors/
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(res.StatusCode, " response. Couldn't read response body: ", err)
		} else {
			log.Println(res.StatusCode, " response: ", string(body))
		}
		return false
	}

	if respBody == nil { // Will be null if the caller doesn't need the response body
		return true
	}

	err = json.NewDecoder(res.Body).Decode(respBody)
	if err != nil {
		log.Println("error while decoding response body: ", err)
		return false
	}
	return true
}

// HttpClient returns an instance of Http Client with custom config
var HttpClient = &http.Client{
	Timeout: time.Second * 20,
}
