package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid CLI config with send enabled",
			config: &Config{
				Mode:             ModeCLI,
				DestEmail:        "test@kindle.com",
				SenderEmail:      "sender@example.com",
				MailjetAPIKey:    "api-key",
				MailjetAPISecret: "api-secret",
				SendEnabled:      true,
			},
			wantErr: false,
		},
		{
			name: "valid CLI config with send disabled",
			config: &Config{
				Mode:        ModeCLI,
				SendEnabled: false,
			},
			wantErr: false,
		},
		{
			name: "valid server config",
			config: &Config{
				Mode:          ModeServer,
				APIKeySecret:  "api-key-secret",
				DynamoDBTable: "test-table",
			},
			wantErr: false,
		},
		{
			name: "server config missing api key",
			config: &Config{
				Mode:          ModeServer,
				DynamoDBTable: "test-table",
			},
			wantErr: true,
		},
		{
			name: "server config missing dynamodb table",
			config: &Config{
				Mode:         ModeServer,
				APIKeySecret: "api-key-secret",
			},
			wantErr: true,
		},
		{
			name: "CLI config missing kindle email with send enabled",
			config: &Config{
				Mode:             ModeCLI,
				SenderEmail:      "sender@example.com",
				MailjetAPIKey:    "api-key",
				MailjetAPISecret: "api-secret",
				SendEnabled:      true,
			},
			wantErr: true,
		},
		{
			name: "CLI config missing sender email with send enabled",
			config: &Config{
				Mode:             ModeCLI,
				DestEmail:        "test@kindle.com",
				MailjetAPIKey:    "api-key",
				MailjetAPISecret: "api-secret",
				SendEnabled:      true,
			},
			wantErr: true,
		},
		{
			name: "CLI config missing mailjet api key with send enabled",
			config: &Config{
				Mode:             ModeCLI,
				DestEmail:        "test@kindle.com",
				SenderEmail:      "sender@example.com",
				MailjetAPISecret: "api-secret",
				SendEnabled:      true,
			},
			wantErr: true,
		},
		{
			name: "CLI config missing mailjet api secret with send enabled",
			config: &Config{
				Mode:          ModeCLI,
				DestEmail:     "test@kindle.com",
				SenderEmail:   "sender@example.com",
				MailjetAPIKey: "api-key",
				SendEnabled:   true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	_ = os.Setenv("F2K_DEST_EMAIL", "test@kindle.com")
	_ = os.Setenv("F2K_SENDER_EMAIL", "sender@example.com")
	_ = os.Setenv("MAILJET_API_KEY", "api-key")
	_ = os.Setenv("MAILJET_API_SECRET", "api-secret")
	_ = os.Setenv("F2K_API_KEY", "api-key-secret")
	_ = os.Setenv("F2K_DYNAMODB_TABLE_NAME", "test-table")
	defer func() {
		_ = os.Unsetenv("F2K_DEST_EMAIL")
		_ = os.Unsetenv("F2K_SENDER_EMAIL")
		_ = os.Unsetenv("MAILJET_API_KEY")
		_ = os.Unsetenv("MAILJET_API_SECRET")
		_ = os.Unsetenv("F2K_API_KEY")
		_ = os.Unsetenv("F2K_DYNAMODB_TABLE_NAME")
	}()

	cfg, err := Load(ModeCLI)
	assert.NoError(t, err)
	assert.Equal(t, "test@kindle.com", cfg.DestEmail)
	assert.Equal(t, "sender@example.com", cfg.SenderEmail)
	assert.Equal(t, "api-key", cfg.MailjetAPIKey)
	assert.Equal(t, "api-secret", cfg.MailjetAPISecret)
	assert.Equal(t, "api-key-secret", cfg.APIKeySecret)
	assert.Equal(t, "test-table", cfg.DynamoDBTable)
	assert.Equal(t, ModeCLI, cfg.Mode)
}

func TestLoadDefaultsToCLI(t *testing.T) {
	cfg, err := Load(ModeCLI)
	assert.NoError(t, err)
	assert.Equal(t, ModeCLI, cfg.Mode)
}

func TestLoadServerMode(t *testing.T) {
	_ = os.Setenv("F2K_API_KEY", "api-key-secret")
	_ = os.Setenv("F2K_DYNAMODB_TABLE_NAME", "test-table")
	defer func() {
		_ = os.Unsetenv("F2K_API_KEY")
		_ = os.Unsetenv("F2K_DYNAMODB_TABLE_NAME")
	}()

	cfg, err := Load(ModeServer)
	assert.NoError(t, err)
	assert.Equal(t, ModeServer, cfg.Mode)
}
