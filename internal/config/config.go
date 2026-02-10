// Package config provides configuration for the savetoink application.
package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/shaftoe/savetoink/internal/constant"
	"github.com/spf13/viper"
)

// Config holds configuration settings for application.
type Config struct {
	DestEmail        string
	SenderEmail      string
	MailjetAPIKey    string
	MailjetAPISecret string
	APIKeySecret     string
	Auth0Domain      string
	Auth0Audience    string
	Debug            bool
	SendEnabled      bool
	DynamoDBTable    string
	Mode             constant.RunMode
	AWSConfig        *aws.Config
	EmailProvider    constant.EmailProvider
	AuthBackend      constant.AuthBackend
}

// Load reads configuration from environment variables and returns a Config instance.
func Load(mode constant.RunMode) (*Config, error) {
	viper.SetEnvPrefix("F2K")
	viper.AutomaticEnv()

	if err := bindEnvVars(); err != nil {
		return nil, err
	}

	cfg := loadConfig(mode)

	if cfg.AuthBackend == "" {
		cfg.AuthBackend = constant.AuthBackendSharedAPIKey
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func bindEnvVars() error {
	envVars := []struct {
		key    string
		envVar string
	}{
		{"api-key", "MAILJET_API_KEY"},
		{"api-key-secret", "F2K_API_KEY"},
		{"api-secret", "MAILJET_API_SECRET"},
		{"auth-backend", "F2K_AUTH_BACKEND"},
		{"auth0-audience", "F2K_AUTH0_AUDIENCE"},
		{"auth0-domain", "F2K_AUTH0_DOMAIN"},
		{"debug", "F2K_DEBUG"},
		{"destination-email", "F2K_DEST_EMAIL"},
		{"dynamodb-table", "F2K_DYNAMODB_TABLE_NAME"},
		{"send-enabled", "F2K_SEND_ENABLED"},
		{"sender-email", "F2K_SENDER_EMAIL"},
	}

	for _, ev := range envVars {
		if err := viper.BindEnv(ev.key, ev.envVar); err != nil {
			return fmt.Errorf("failed to bind %s env: %w", ev.key, err)
		}
	}
	return nil
}

func loadConfig(mode constant.RunMode) *Config {
	cfg := &Config{
		DestEmail:        viper.GetString("destination-email"),
		SenderEmail:      viper.GetString("sender-email"),
		MailjetAPIKey:    viper.GetString("api-key"),
		MailjetAPISecret: viper.GetString("api-secret"),
		APIKeySecret:     viper.GetString("api-key-secret"),
		Auth0Domain:      viper.GetString("auth0-domain"),
		Auth0Audience:    viper.GetString("auth0-audience"),
		Debug:            viper.GetBool("debug"),
		SendEnabled:      viper.GetBool("send-enabled"),
		DynamoDBTable:    viper.GetString("dynamodb-table"),
		Mode:             mode,
		AuthBackend:      constant.AuthBackend(viper.GetString("auth-backend")),
	}

	_ = viper.BindEnv("api-key", "MJ_APIKEY_PUBLIC")
	_ = viper.BindEnv("api-secret", "MJ_APIKEY_PRIVATE")

	return cfg
}

// Validate checks that all required configuration fields are set.
func (c *Config) Validate() error {
	var missing []string

	if c.Mode == constant.ModeServer {
		if err := c.validateServerConfig(&missing); err != nil {
			return err
		}
	}

	if c.SendEnabled {
		c.EmailProvider = constant.EmailBackendMailjet
		c.validateSendEnabledConfig(&missing)
	}

	if len(missing) > 0 {
		return fmt.Errorf("required environment variables are missing: %v", missing)
	}

	return nil
}

func (c *Config) validateServerConfig(missing *[]string) error {
	switch c.AuthBackend {
	case constant.AuthBackendSharedAPIKey:
		if c.APIKeySecret == "" {
			*missing = append(*missing, "F2K_API_KEY")
		}
	case constant.AuthBackendAuth0:
		if c.Auth0Domain == "" {
			*missing = append(*missing, "F2K_AUTH0_DOMAIN")
		}
		if c.Auth0Audience == "" {
			*missing = append(*missing, "F2K_AUTH0_AUDIENCE")
		}
	default:
		return fmt.Errorf("unsupported auth backend: %s", c.AuthBackend)
	}
	if c.DynamoDBTable == "" {
		*missing = append(*missing, "F2K_DYNAMODB_TABLE_NAME")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	c.AWSConfig = &cfg
	return nil
}

func (c *Config) validateSendEnabledConfig(missing *[]string) {
	if c.DestEmail == "" {
		*missing = append(*missing, "F2K_DEST_EMAIL")
	}
	if c.SenderEmail == "" {
		*missing = append(*missing, "F2K_SENDER_EMAIL")
	}
	if c.MailjetAPIKey == "" {
		*missing = append(*missing, "MAILJET_API_KEY")
	}
	if c.MailjetAPISecret == "" {
		*missing = append(*missing, "MAILJET_API_SECRET")
	}
}
