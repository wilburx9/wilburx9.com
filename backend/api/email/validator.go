package email

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
)

func validateData(data Data) string {
	message := ""
	_, err := mail.ParseAddress(data.SenderEmail)
	switch {
	case err != nil:
		message = "Please, enter a valid email address"
	case len(strings.TrimSpace(data.SenderName)) == 0:
		message = "Please, enter your name"
	case len(strings.TrimSpace(data.Subject)) == 0:
		message = "Please, enter a subject"
	case len(strings.TrimSpace(data.Message)) == 0:
		message = "Please, enter a message"
	case len(strings.TrimSpace(data.CaptchaResponse)) == 0:
		message = "Cannot verify humanness"
	}
	return message
}

func validateRecaptcha(response string, httpClient internal.HttpClient) bool {
	u := "https://hcaptcha.com/siteverify"
	payload := url.Values{}
	payload.Add("response", strings.TrimSpace(response))
	payload.Add("secret", configs.Config.HCaptchaSecret)
	payload.Add("sitekey", configs.Config.HCaptchaSiteKey)

	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(payload.Encode()))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return false
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return false
	}
	defer res.Body.Close()

	resp := struct {
		Success bool `json:"success"`
	}{}
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return false
	}
	return resp.Success
}
