package configs

import (
	"fmt"
	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config is global object that holds all application level variables.
var Config appConfig

type appConfig struct {
	Port                 string `mapstructure:"port"`
	MediumUsername       string `mapstructure:"medium_username"`
	WPUrl                string `mapstructure:"wp_url"`
	UnsplashUsername     string `mapstructure:"unsplash_username"`
	UnsplashAccessKey    string `mapstructure:"unsplash_access_key"`
	InstagramAccessToken string `mapstructure:"instagram_access_token"`
	Env                  string `mapstructure:"env"`
	SentryDsn            string `mapstructure:"sentry_dsn"`
	GithubToken          string `mapstructure:"github_token"`
	GithubUsername       string `mapstructure:"github_username"`
	SmtpHost             string `mapstructure:"smtp_host"`
	SmtpPort             int    `mapstructure:"smtp_port"`
	SmtpUsername         string `mapstructure:"smtp_username"`
	SmtpPassword         string `mapstructure:"smtp_password"`
	ContactEmail         string `mapstructure:"contact_email"`
	RecaptchaSecret      string `mapstructure:"recaptcha_secret"`
}

// IsRelease returns true for release Env and false otherwise
func (c appConfig) IsRelease() bool {
	return c.Env == "release"
}

// IsDebug returns true for debug Env and false otherwise
func (c appConfig) IsDebug() bool {
	return c.Env == "debug"
}

// LoadConfig loads config variables from a config file or environment variables
func LoadConfig(path string) error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("WilburX9")
	v.AutomaticEnv()
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}

	err := v.Unmarshal(&Config)
	log.WithFields(structs.Map(Config)).Info("App started with these config")
	return err
}
