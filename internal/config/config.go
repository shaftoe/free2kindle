// Package config provides configuration for the free2kindle application.
package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"
)

// RunMode defines the application execution mode.
type RunMode string

const (
	// ModeCLI indicates CLI execution mode.
	ModeCLI RunMode = "cli"
	// ModeServer indicates server execution mode.
	ModeServer RunMode = "server"
)

// Config holds configuration settings for application.
type Config struct {
	Account          string
	DestEmail        string
	SenderEmail      string
	MailjetAPIKey    string
	MailjetAPISecret string
	APIKeySecret     string
	Debug            bool
	SendEnabled      bool
	DynamoDBTable    string
	Mode             RunMode
	AWSConfig        *aws.Config
	EmailProvider    string
}

// Load reads configuration from environment variables and returns a Config instance.
func Load(mode RunMode) (*Config, error) {
	viper.SetEnvPrefix("F2K")
	viper.AutomaticEnv()

	if err := viper.BindEnv("destination-email", "F2K_DEST_EMAIL"); err != nil {
		return nil, fmt.Errorf("failed to bind destination-email env: %w", err)
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
	if err := viper.BindEnv("debug", "F2K_DEBUG"); err != nil {
		return nil, fmt.Errorf("failed to bind debug env: %w", err)
	}
	if err := viper.BindEnv("send-enabled", "F2K_SEND_ENABLED"); err != nil {
		return nil, fmt.Errorf("failed to bind send-enabled env: %w", err)
	}
	if err := viper.BindEnv("dynamodb-table", "F2K_DYNAMODB_TABLE_NAME"); err != nil {
		return nil, fmt.Errorf("failed to bind dynamodb-table env: %w", err)
	}

	cfg := &Config{
		Account:          "free2kindle", // currently hardcoded, auth work in progress
		DestEmail:        viper.GetString("destination-email"),
		SenderEmail:      viper.GetString("sender-email"),
		MailjetAPIKey:    viper.GetString("api-key"),
		MailjetAPISecret: viper.GetString("api-secret"),
		APIKeySecret:     viper.GetString("api-key-secret"),
		Debug:            viper.GetBool("debug"),
		SendEnabled:      viper.GetBool("send-enabled"),
		DynamoDBTable:    viper.GetString("dynamodb-table"),
		Mode:             mode,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required configuration fields are set.
func (c *Config) Validate() error {
	var missing []string

	if c.Mode == ModeServer {
		if c.APIKeySecret == "" {
			missing = append(missing, "F2K_API_KEY")
		}
		if c.DynamoDBTable == "" {
			missing = append(missing, "F2K_DYNAMODB_TABLE_NAME")
		}

		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}
		c.AWSConfig = &cfg
	}

	if c.SendEnabled {
		if c.DestEmail == "" {
			missing = append(missing, "F2K_DEST_EMAIL")
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
		c.EmailProvider = "MailJet"
	}

	if len(missing) > 0 {
		return fmt.Errorf("required environment variables are missing: %v", missing)
	}

	return nil
}
