package contact

import (
	"net/mail"
	"strings"
)

func validate(data requestData) string {
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
	}
	return message
}
