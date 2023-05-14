package common

import (
	"context"
	"os"
)

// Config is a container of environment variables
type Config struct {
	TurnstileSecret    string
	TurnstileHostName  string
	MailChimpToken     string
	MailChimpDC        string
	NewsletterListId   string
	EmailSender        string
	ProgrammingSegment string
	PhotographySegment string
}

// NewConfig instantiates Config
func NewConfig() *Config {
	return &Config{
		TurnstileSecret:    os.Getenv("TURNSTILE_SECRET"),
		TurnstileHostName:  os.Getenv("TURNSTILE_HOSTNAME"),
		MailChimpToken:     os.Getenv("MAILCHIMP_TOKEN"),
		MailChimpDC:        os.Getenv("MAILCHIMP_DATA_CENTRE"),
		NewsletterListId:   os.Getenv("NEWSLETTER_LIST_ID"),
		EmailSender:        os.Getenv("EMAIL_SENDER"),
		ProgrammingSegment: os.Getenv("PROGRAMMING_SEGMENT"),
		PhotographySegment: os.Getenv("PHOTOGRAPHY_SEGMENT"),
	}
}

func ConfigFromContext(ctx *context.Context) *Config {
	return (*ctx).Value(ConfigKey).(*Config)
}

// ConfigKey is the key for the config in the context
var ConfigKey = "common.Config"
