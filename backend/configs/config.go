package configs

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
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
	GithubToken          string `mapstructure:"github_token"`
	GithubUsername       string `mapstructure:"github_username"`
	AppHome              string `mapstructure:"app_home"`
	GcpProjectId         string `mapstructure:"gcp_project_id"`
	EmailDomain          string `mapstructure:"email_domain"`
	EmailReceiver        string `mapstructure:"email_receiver"`
	EmailAPIKey          string `mapstructure:"email_api_key"`
	HCaptchaSecret       string `mapstructure:"h_captcha_secret"`
	HCaptchaSiteKey      string `mapstructure:"h_captcha_site_key"`
	APIKey               string `mapstructure:"api_key"`
	APISalt              string `mapstructure:"api_salt"`
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
func LoadConfig() error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("WilburX9")
	v.AddConfigPath(fmt.Sprintf("%vconfigs", os.Getenv("WILBURX9_APP_HOME")))
	v.AutomaticEnv()

	if err := v.BindEnv("port", "PORT"); err != nil {
		return fmt.Errorf("unable to read PORT from env: %s", err)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}

	err := v.Unmarshal(&Config)
	log.Infof("App started with these config: %+v\n", Config)
	return err
}
