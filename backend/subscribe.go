package main

import (
	"backend/common"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mailerlite/mailerlite-go"
	"log"
	"net/http"
	"net/mail"
	"strings"
)

func main() {
	lambda.Start(handleSubscribe)
}

// handleSubscribe is called when the Lambda receivers a request.
// nil errors are returned because I want return custom http errors as opposed to Lambda's default 500.
func handleSubscribe(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	status, msg := processSubscribeRequest(ctx, req.Body)
	origin := req.Headers["origin"]

	return common.MakeResponse(origin, status, msg), nil
}

func processSubscribeRequest(ctx context.Context, body string) (int, string) {
	data, msg, err := validateForm(body)
	if msg != "" || err != nil {
		log.Println("Failed to validate request body",
			"error: ", fmt.Sprintf("%v", err),
			"message: ", msg)
		return http.StatusBadRequest, msg
	}

	err = validateCaptcha(ctx, data.Captcha)
	if err != nil {
		log.Println("Failed to validate captcha",
			"error: ", err.Error(),
		)
		return http.StatusBadRequest, "Unable to complete subscription"
	}

	err = subscribe(ctx, data.Email, data.Tags)
	if err != nil {
		log.Println("Subscription request failed",
			"error: ", err.Error(),
		)
		return http.StatusBadGateway, "Something went wrong"
	}

	return http.StatusCreated, "Successfully created"
}

// subscribe forwards the request to the mail client to subscribe the user
func subscribe(ctx context.Context, email string, tags []string) error {

	allGroups, _, err := common.MailClient.Group.List(ctx, &mailerlite.ListGroupOptions{})
	if err != nil {
		return err
	}

	var groups []string
	// Filter out groups with same name as supported tags
	for _, group := range allGroups.Data {
		for _, tag := range tags {
			if strings.EqualFold(group.Name, tag) {
				groups = append(groups, group.ID)
			}
		}
	}

	u := "https://connect.mailerlite.com/api/subscribers"
	subscriber := map[string]interface{}{"email": email, "groups": groups, "status": "unconfirmed"}
	reqBody := new(bytes.Buffer)

	err = json.NewEncoder(reqBody).Encode(subscriber)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", common.MailClient.APIKey()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if err == nil && (res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
		return nil
	}

	return fmt.Errorf("subscription request failed. status code: %v, error: %v", res.StatusCode, err)
}

// validateCaptcha ensures this is not a spam request
func validateCaptcha(ctx context.Context, captcha string) error {
	data := fmt.Sprintf("secret=%v&response=%v", common.AppConfig.TurnstileSecret, captcha)

	u := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := common.HttpClient.Do(req)

	defer res.Body.Close()

	var t turnstile
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return err
	}

	if t.Success && t.Hostname == common.AppConfig.TurnstileHostName {
		return nil
	}

	return fmt.Errorf("invalid response: %+v\n", t)

}

// validateForm confirms that the request body contains the required data in the valid formats.
func validateForm(body string) (requestData, string, error) {
	var data requestData
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return requestData{}, "invalid request body", err
	}

	address, err := mail.ParseAddress(data.Email)
	if err != nil {
		return requestData{}, "invalid email", err
	}

	if strings.TrimSpace(data.Captcha) == "" {
		return requestData{}, "captcha is required", nil
	}

	if len(data.Tags) > 2 {
		return requestData{}, "invalid tags", nil
	}

	data.Email = address.Address
	data.Tags = cleanTags(data.Tags)

	return data, "", nil
}

// cleanTags ensures the tags in the request are valid. Also, adds default tags
func cleanTags(rawTags []string) []string {
	var tapsMap = make(map[string]bool, 0) // Use a map to prevent duplicates

	for _, tag := range rawTags {
		trimmed := strings.ToLower(strings.TrimSpace(tag))
		// Only take the tag if it's valid
		if trimmed == common.Photography || trimmed == common.Programming {
			tapsMap[trimmed] = true
		}
	}

	var tags = []string{common.Blog} // Every subscriber belongs to the "blog" tag

	// Convert the map to a list
	for v := range tapsMap {
		tags = append(tags, v)
	}

	// If tags wasn't sent in the request, add all supported tags
	if len(tags) == 1 {
		tags = append(tags, common.Photography, common.Programming)
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
