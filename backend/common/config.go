package common

import (
	"os"
)

// Config is a container of environment variables
type Config struct {
	TurnstileSecret   string
	TurnstileHostName string
	NewsletterListId  string
	EmailSender       string
	AllowedOrigins    []string
	MailerLiteToken   string
	TimeZone          string
}

// NewConfig instantiates Config
func NewConfig() *Config {
	return &Config{
		TurnstileSecret:   os.Getenv("TURNSTILE_SECRET"),
		TurnstileHostName: os.Getenv("TURNSTILE_HOSTNAME"),
		NewsletterListId:  os.Getenv("NEWSLETTER_LIST_ID"),
		EmailSender:       os.Getenv("EMAIL_SENDER"),
		MailerLiteToken:   os.Getenv("MAILER_LITE_TOKEN"),
		TimeZone:          os.Getenv("TZ"),
		AllowedOrigins:    []string{os.Getenv("PROD_FRONTEND_URL"), os.Getenv("DEV_FRONTEND_URL")},
	}
}
