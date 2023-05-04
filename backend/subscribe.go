package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wilburt/wilburx9.com/backend"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

const (
	photography = "photography"
	programming = "programming"
)

func main() {
	lambda.Start(start)
}

func start(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data, err := validateForm(req.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       GetResponseBody(false, err.Error()),
		}, nil
	}

	err = validateCaptcha(data.captcha)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Body:       GetResponseBody(false, "Unable to complete subscription"),
		}, nil
	}

	err = subscribe(data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body:       GetResponseBody(false, "Something went wrong"),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       GetResponseBody(true, "Successfully created"),
	}, nil
}

func subscribe(data formData) error {
	return errors.New("Not implemented yet")
}

func validateCaptcha(captcha string) error {
	data := fmt.Sprintf("secret=%v&response=%v", os.Getenv(""), captcha)

	u := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(req)

	defer res.Body.Close()

	var t turnstile
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return err
	}

	if t.Success && t.Hostname == "wilburx9.com" {
		return nil
	}

	return fmt.Errorf("invalid response: %+v\\n", t)

}

func validateForm(body string) (formData, error) {
	form, err := url.ParseQuery(body)
	if err != nil {
		return formData{}, err
	}

	email := form.Get("email")
	rawTags := form["tags"]
	captcha := form.Get("captcha")

	if strings.TrimSpace(email) == "" {
		return formData{}, errors.New("email is required")
	}

	if strings.TrimSpace(captcha) == "" {
		return formData{}, errors.New("captcha is required")
	}

	return formData{
		email:   email,
		captcha: captcha,
		tags:    getTags(rawTags),
	}, nil
}

func getTags(rawTags []string) []string {
	sort.Strings(rawTags)
	var tags = make([]string, 2)
	if strings.ToLower(strings.TrimSpace(rawTags[0])) == photography {
		tags = append(tags, photography)
	}

	if strings.ToLower(strings.TrimSpace(rawTags[1])) == programming {
		tags = append(tags, programming)
	}

	if len(tags) == 0 {
		tags = append(tags, photography, programming)
	}

	return tags

}

type formData struct {
	email   string
	captcha string
	tags    []string
}

type turnstile struct {
	Success    bool          `json:"success"`
	Hostname   string        `json:"hostname"`
	ErrorCodes []interface{} `json:"error-codes"` // TODO: Remove
	Action     string        `json:"action"`
}
