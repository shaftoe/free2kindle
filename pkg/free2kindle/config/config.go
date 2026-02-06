// Package config provides configuration for the free2kindle application.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the configuration settings for the application.
type Config struct {
	KindleEmail      string
	SenderEmail      string
	MailjetAPIKey    string
	MailjetAPISecret string
	APIKeySecret     string
	Debug            bool
	SendEnabled      bool
}

// Load reads configuration from environment variables and returns a Config instance.
func Load() (*Config, error) {
	viper.SetEnvPrefix("F2K")
	viper.AutomaticEnv()

	if err := viper.BindEnv("kindle-email", "F2K_KINDLE_EMAIL"); err != nil {
		return nil, fmt.Errorf("failed to bind kindle-email env: %w", err)
	}
	if err := viper.BindEnv("sender-email", "F2K_SENDER_EMAIL"); err != nil {
		return nil, fmt.Errorf("failed to bind sender-email env: %w", err)
	}
	if err := viper.BindEnv("api-key", "MAILJET_API_KEY", "MJ_APIKEY_PUBLIC"); err != nil {
		return nil, fmt.Errorf("failed to bind api-key env: %w", err)
	}
	if err := viper.BindEnv("api-secret", "MAILJET_API_SECRET", "MJ_APIKEY_PRIVATE"); err != nil {
		return nil, fmt.Errorf("failed to bind api-secret env: %w", err)
	}
	if err := viper.BindEnv("api-key-secret", "F2K_API_KEY"); err != nil {
		return nil, fmt.Errorf("failed to bind api-key-secret env: %w", err)
	}
	if err := viper.BindEnv("debug", "DEBUG"); err != nil {
		return nil, fmt.Errorf("failed to bind debug env: %w", err)
	}
	if err := viper.BindEnv("send-enabled", "F2K_SEND_ENABLED"); err != nil {
		return nil, fmt.Errorf("failed to bind send-enabled env: %w", err)
	}

	cfg := &Config{
		KindleEmail:      viper.GetString("kindle-email"),
		SenderEmail:      viper.GetString("sender-email"),
		MailjetAPIKey:    viper.GetString("api-key"),
		MailjetAPISecret: viper.GetString("api-secret"),
		APIKeySecret:     viper.GetString("api-key-secret"),
		Debug:            viper.GetBool("debug"),
		SendEnabled:      viper.GetBool("send-enabled"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required configuration fields are set.
func (c *Config) Validate() error {
	var missing []string

	if c.APIKeySecret == "" {
		missing = append(missing, "F2K_API_KEY")
	}

	if c.SendEnabled {
		if c.KindleEmail == "" {
			missing = append(missing, "F2K_KINDLE_EMAIL")
		}
		if c.SenderEmail == "" {
			missing = append(missing, "F2K_SENDER_EMAIL")
		}
		if c.MailjetAPIKey == "" {
			missing = append(missing, "MAILJET_API_KEY")
		}
		if c.MailjetAPISecret == "" {
			missing = append(missing, "MAILJET_API_SECRET")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required environment variables are missing: %v", missing)
	}

	return nil
}
