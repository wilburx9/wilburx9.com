package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	. "github.com/wilburt/wilburx9.com/backend/common"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
)

const (
	photography = "photography"
	programming = "programming"
	blog        = "blog"
)

func main() {
	lambda.Start(start)
}

// start is called when the Lambda receivers a request
func start(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data, msg := validateForm(req.Body)
	if msg != "" {
		return MakeResponse(http.StatusBadRequest, msg), nil
	}

	err := validateCaptcha(data.Captcha)
	if err != nil {
		log.Println(err)
		return MakeResponse(http.StatusUnprocessableEntity, "Unable to complete subscription"), nil
	}

	err = subscribe(data)
	if err != nil {
		log.Println(err)
		return MakeResponse(http.StatusBadGateway, "Something went wrong"), nil
	}

	return MakeResponse(
		http.StatusCreated,
		"Successfully created",
	), nil
}

// subscribe forwards the request to MailChimp to subscribe the user
func subscribe(data requestData) error {
	log.Printf("tryig to subscribe with %+v\n", data)
	dc := os.Getenv("MAILCHIMP_DC")
	token := os.Getenv("MAILCHIMP_TOKEN")
	listId := os.Getenv("MAILCHIMP_LIST_ID")
	member := map[string]interface{}{"email_address": data.Email, "status": "pending", "tags": data.Tags}
	u := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members", dc, listId)

	log.Printf("tryig to subscribe with member %+v\n", member)

	reqBody, err := json.Marshal(member)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if err == nil && res.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("subscription request returneed an error or non-200 status code. %v :: %w", res.StatusCode, err)
}

// validateCaptcha ensures this is not a spam request
func validateCaptcha(captcha string) error {
	secret := os.Getenv("TURNSTILE_SECRET")
	data := fmt.Sprintf("secret=%v&response=%v", secret, captcha)

	u := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := HttpClient.Do(req)

	defer res.Body.Close()

	var t turnstile
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return err
	}

	if t.Success && t.Hostname == "wilburx9.com" {
		return nil
	}

	return fmt.Errorf("invalid response: %+v\n", t)

}

// validateForm confirms that the request body contains the required data in the valid formats.
func validateForm(body string) (requestData, string) {
	var data requestData
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return requestData{}, "Bad Request"
	}

	address, err := mail.ParseAddress(data.Email)
	if err != nil {
		return requestData{}, "Invalid email"
	}

	if strings.TrimSpace(data.Captcha) == "" {
		return requestData{}, "Captcha is required"
	}

	if len(data.Tags) > 2 {
		return requestData{}, "Invalid tags"
	}

	data.Email = address.Address
	data.Tags = cleanTags(data.Tags)

	return data, ""
}

// cleanTags ensures the tags in the request are valid. Also, adds default tags
func cleanTags(rawTags []string) []string {
	var tapsMap = make(map[string]bool, 0) // Use a map to prevent duplicates

	for _, tag := range rawTags {
		trimmed := strings.ToLower(strings.TrimSpace(tag))
		// Only take the tag if it's valid
		if trimmed == photography || trimmed == programming {
			tapsMap[trimmed] = true
		}
	}

	var tags = []string{blog} // Every subscriber belongs to the "blog" tag

	// Convert the map to a list
	for v := range tapsMap {
		tags = append(tags, v)
	}

	// If tags wasn't sent in the request, add all supported tags
	if len(tags) == 1 {
		tags = append(tags, photography, programming)
	}

	return tags

}

type requestData struct {
	Email   string   `json:"email"`
	Captcha string   `json:"captcha"`
	Tags    []string `json:"tags"`
}

type turnstile struct {
	Success  bool   `json:"success"`
	Hostname string `json:"hostname"`
	Action   string `json:"action"`
}
