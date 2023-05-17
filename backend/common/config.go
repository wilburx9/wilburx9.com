package common

import (
	"os"
	"strings"
)

// Config is a container of environment variables
type Config struct {
	TurnstileSecret   string
	TurnstileHostName string
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
		EmailSender:       os.Getenv("EMAIL_SENDER"),
		MailerLiteToken:   os.Getenv("MAILER_LITE_TOKEN"),
		TimeZone:          strings.ReplaceAll(os.Getenv("TZ"), ":", ""),
		AllowedOrigins:    []string{os.Getenv("PROD_FRONTEND_URL"), os.Getenv("DEV_FRONTEND_URL")},
	}
}
