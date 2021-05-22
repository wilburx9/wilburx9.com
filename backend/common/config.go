package common

import (
	"fmt"
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
	GcpSaKey             string `mapstructure:"gcp_backend_sa_key"`
	GcpProjectId         string `mapstructure:"gcp_project_id"`
}

func (c appConfig) isRelease() bool {
	return c.Env == "release"
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

	return v.Unmarshal(&Config)
}
