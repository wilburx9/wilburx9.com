package common

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"os"
	"reflect"
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

// newConfig instantiates Config
func newConfig() (*Config, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("aws session init failed: %w", err)
	}

	ssmSvc := ssm.New(sess)

	m := map[string]any{
		"allowed_origins":    []string{},
		"email_sender":       "",
		"mailer_lite_token":  "",
		"turnstile_hostname": "",
		"turnstile_secret":   "",
	}

	// Read all the secrets into the map
	for k, v := range m {
		key := strings.ToUpper(fmt.Sprintf("wilburx9_%v", k))
		input := &ssm.GetParameterInput{
			Name:           aws.String(key),
			WithDecryption: aws.Bool(true),
		}
		param, err := ssmSvc.GetParameter(input)
		if err != nil {
			return nil, fmt.Errorf("unable to read %q from store: %w in %v", key, err, *sess.Config.Region)
		}
		if (reflect.TypeOf(v)).Kind() == reflect.Slice {
			m[k] = strings.Split(*param.Parameter.Value, ",")
		} else {
			m[k] = *param.Parameter.Value
		}
	}

	mBytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal config map: %w", err)
	}

	var config Config
	err = json.Unmarshal(mBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal config map: %w", err)
	}

	config.TimeZone = strings.ReplaceAll(os.Getenv("TZ"), ":", "")
	return &config, nil
}
