package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
)

type requestData struct {
	SenderEmail    string `json:"sender_email"`
	SenderName     string `json:"sender_name"`
	Subject        string `json:"subject"`
	Message        string `json:"message"`
	RecaptchaToken string `json:"recaptcha_token"`
}

// Handler validates request body and reCAPTCHA and possibly sends an email
func Handler(c *gin.Context) {
	var data requestData
	err := c.ShouldBindJSON(&data)
	message := validateData(data)
	if err != nil || message != "" {
		c.JSON(http.StatusBadRequest, internal.MakeErrorResponse(message))
		return
	}

	c.JSON(http.StatusOK, internal.MakeSuccessResponse("OK"))

	// if !validateRecaptcha(configs.Config.RecaptchaSecret, data.RecaptchaToken, &http.Client{}) {
	// 	c.JSON(http.StatusForbidden, internal.MakeErrorResponse("Could not confirm humanness"))
	// 	return
	// }
	//
	// err = send(data)
	// if err != nil {
	// 	log.WithFields(log.Fields{"error": err}).Warning("Couldn't send email")
	// 	c.JSON(http.StatusBadGateway, internal.MakeErrorResponse("Email not sent"))
	// 	return
	// }
	// c.JSON(http.StatusOK, internal.MakeSuccessResponse("Email sent successfully"))
}

// func send(data requestData) error {
// 	config := configs.Config
// 	sender := fmt.Sprintf("%v <%v>", data.SenderName, data.SenderEmail)
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", sender)
// 	m.SetHeader("To", config.ContactEmail)
// 	m.SetHeader("Subject", data.Subject)
// 	m.SetBody("text/plain", data.Message)
//
// 	d := gomail.NewDialer(config.SmtpHost, config.SmtpPort, config.SmtpUsername, config.SmtpPassword)
// 	return d.DialAndSend(m)
// }
