package contact

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
	"net/mail"
	"strings"
)

const (
	minRecaptchaThresh = 0.5
)

func validateData(data requestData) string {
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
	case len(strings.TrimSpace(data.RecaptchaToken)) == 0:
		message = "Cannot verify humanness"
	}
	return message
}

func validateRecaptcha(secret string, token string, httpClient internal.HttpClient) bool {
	url := fmt.Sprintf("https://www.google.com/recaptcha/api/siteverify?secret=%v&response=%v", secret, token)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return false
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return false
	}
	defer res.Body.Close()

	resp := struct {
		Success bool    `json:"success"`
		Score   float64 `json:"score"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
		return false
	}

	return resp.Success && resp.Score >= minRecaptchaThresh
}
