package main

import (
	"backend/common"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
)

func main() {
	lambda.Start(handleSubscribe)
}

// handleSubscribe is called when the Lambda receivers a request.
// nil errors are returned because I want return custom http errors as opposed to Lambda's default 500.
func handleSubscribe(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data, msg, err := validateForm(req.Body)
	if msg != "" || err != nil {
		log.Println("Failed to validate request body",
			"error: ", fmt.Sprintf("%v", err),
			"message: ", msg)
		return common.MakeResponse(http.StatusBadRequest, msg), nil
	}

	err = validateCaptcha(ctx, data.Captcha)
	if err != nil {
		log.Println("Failed to validate captcha",
			"error: ", err.Error(),
		)
		return common.MakeResponse(http.StatusUnprocessableEntity, "Unable to complete subscription"), nil
	}

	err = subscribe(ctx, data)
	if err != nil {
		log.Println("Subscription request failed",
			"error: ", err.Error(),
		)
		return common.MakeResponse(http.StatusBadGateway, "Something went wrong"), nil
	}

	return common.MakeResponse(
		http.StatusCreated,
		"Successfully created",
	), nil
}

// subscribe forwards the request to MailChimp to subscribe the user
func subscribe(ctx context.Context, data requestData) error {
	listId := os.Getenv("MAILCHIMP_LIST_ID")
	member := map[string]interface{}{"email_address": data.Email, "status": "pending", "tags": data.Tags}

	err := common.MakeMailChimpRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("lists/%s/members", listId),
		member,
		nil,
	)
	if err != nil {
		return fmt.Errorf("subscription request returneed an error: %w", err)
	}
	return nil
}

// validateCaptcha ensures this is not a spam request
func validateCaptcha(ctx context.Context, captcha string) error {
	secret := os.Getenv("TURNSTILE_SECRET")
	data := fmt.Sprintf("secret=%v&response=%v", secret, captcha)

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

	if t.Success && t.Hostname == os.Getenv("TURNSTILE_HOSTNAME") {
		return nil
	}

	return fmt.Errorf("invalid response: %+v\n", t)

}

// validateForm confirms that the request body contains the required data in the valid formats.
func validateForm(body string) (requestData, string, error) {
	var data requestData
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return requestData{}, "Bad Request", err
	}

	address, err := mail.ParseAddress(data.Email)
	if err != nil {
		return requestData{}, "Invalid email", err
	}

	if strings.TrimSpace(data.Captcha) == "" {
		return requestData{}, "Captcha is required", nil
	}

	if len(data.Tags) > 2 {
		return requestData{}, "Invalid tags", nil
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
