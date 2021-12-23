package contact

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/configs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type requestData struct {
	SenderEmail    string `json:"sender_email"`
	SenderName     string `json:"sender_name"`
	Subject        string `json:"subject"`
	Message        string `json:"message"`
	RecaptchaToken string `json:"recaptcha_token"`
}

// Handler validates request body and reCAPTCHA and possibly sends an email
func Handler(c *gin.Context, client internal.HttpClient) {
	var data requestData
	err := c.ShouldBindJSON(&data)
	message := validateData(data)
	if err != nil || message != "" {
		c.JSON(http.StatusBadRequest, internal.MakeErrorResponse(message))
		return
	}

	// if !validateRecaptcha(configs.Config.RecaptchaSecret, data.RecaptchaToken, &http.Client{}) {
	// 	c.JSON(http.StatusForbidden, internal.MakeErrorResponse("Could not confirm humanness"))
	// 	return
	// }

	err = send(data, client)
	if err != nil {
		c.JSON(http.StatusBadGateway, internal.MakeErrorResponse("Email not sent"))
		return
	}

	c.JSON(http.StatusOK, internal.MakeSuccessResponse("Email sent successfully"))
}

func send(data requestData, client internal.HttpClient) error {
	u := fmt.Sprintf("https://api.mailgun.net/v3/%v/messages", configs.Config.EmailDomain)
	payload := url.Values{}
	payload.Add("from", fmt.Sprintf("%v <%v>", data.SenderName, data.SenderEmail))
	payload.Add("to", configs.Config.EmailReceiver)
	payload.Add("subject", data.Subject)
	payload.Add("text", data.Message)

	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(payload.Encode()))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("api:%v", configs.Config.EmailAPIKey)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", auth))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return err
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fields := log.Fields{"data": data, "statusCode": res.StatusCode, "response": string(body)}
	log.WithFields(fields).Info("Post email")

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("could not send email")
	}
	return nil
}
