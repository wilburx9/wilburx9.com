package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"log"
	"os"
	"strings"
)

// Config is a container of environment variables
type Config struct {
	TurnstileSecret   string   `json:"turnstile_secret"`
	TurnstileHostname string   `json:"turnstile_hostname"`
	EmailSender       string   `json:"email_sender"`
	AllowedOrigins    []string `json:"allowed_origins"`
	MailerLiteToken   string   `json:"mailer_lite_token"`
	TimeZone          string   `json:"time_zone"`
}

// InitConfig instantiates Config
func InitConfig() error {
	config, err := newConfig()
	if err != nil {
		log.Println(err)
		return errors.New("something went wrong")
	}

	AppConfig = config
	return nil
}

func newConfig() (*Config, error) {
	sess, err := session.NewSession()
	if err != nil {
		return &Config{}, fmt.Errorf("aws session init failed: %w", err)
	}

	svc := secretsmanager.New(sess)
	m := map[string]any{
		"allowed_origins":    []string{},
		"email_sender":       "",
		"mailer_lite_token":  "",
		"turnstile_hostname": "",
		"turnstile_secret":   "",
	}

	for s := range m {
		key := strings.ToUpper(fmt.Sprintf("wilburx9_%v", s))
		input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(key)}
		value, err := svc.GetSecretValue(input)
		if err != nil {
			return &Config{}, fmt.Errorf("unable to read %q from store: %w", key, err)
		}
		m[s] = value
	}

	mBytes, err := json.Marshal(m)
	if err != nil {
		return &Config{}, fmt.Errorf("unable to marshal config map: %w", err)
	}

	var config Config
	err = json.Unmarshal(mBytes, &config)
	if err != nil {
		return &Config{}, fmt.Errorf("unable to marshal config map: %w", err)
	}

	config.TimeZone = strings.ReplaceAll(os.Getenv("TZ"), ":", "")
	return &config, nil
}
