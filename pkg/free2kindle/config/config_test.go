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
			name: "valid config",
			config: &Config{
				KindleEmail:      "test@kindle.com",
				SenderEmail:      "sender@example.com",
				MailjetAPIKey:    "api-key",
				MailjetAPISecret: "api-secret",
				APIKeySecret:     "api-key-secret",
			},
			wantErr: false,
		},
		{
			name: "missing kindle email",
			config: &Config{
				SenderEmail:      "sender@example.com",
				MailjetAPIKey:    "api-key",
				MailjetAPISecret: "api-secret",
				APIKeySecret:     "api-key-secret",
			},
			wantErr: true,
		},
		{
			name: "missing sender email",
			config: &Config{
				KindleEmail:      "test@kindle.com",
				MailjetAPIKey:    "api-key",
				MailjetAPISecret: "api-secret",
				APIKeySecret:     "api-key-secret",
			},
			wantErr: true,
		},
		{
			name: "missing mailjet api key",
			config: &Config{
				KindleEmail:      "test@kindle.com",
				SenderEmail:      "sender@example.com",
				MailjetAPISecret: "api-secret",
				APIKeySecret:     "api-key-secret",
			},
			wantErr: true,
		},
		{
			name: "missing mailjet api secret",
			config: &Config{
				KindleEmail:   "test@kindle.com",
				SenderEmail:   "sender@example.com",
				MailjetAPIKey: "api-key",
				APIKeySecret:  "api-key-secret",
			},
			wantErr: true,
		},
		{
			name: "missing api key secret",
			config: &Config{
				KindleEmail:      "test@kindle.com",
				SenderEmail:      "sender@example.com",
				MailjetAPIKey:    "api-key",
				MailjetAPISecret: "api-secret",
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
	_ = os.Setenv("F2K_KINDLE_EMAIL", "test@kindle.com")
	_ = os.Setenv("F2K_SENDER_EMAIL", "sender@example.com")
	_ = os.Setenv("MAILJET_API_KEY", "api-key")
	_ = os.Setenv("MAILJET_API_SECRET", "api-secret")
	_ = os.Setenv("F2K_API_KEY", "api-key-secret")
	defer func() {
		_ = os.Unsetenv("F2K_KINDLE_EMAIL")
		_ = os.Unsetenv("F2K_SENDER_EMAIL")
		_ = os.Unsetenv("MAILJET_API_KEY")
		_ = os.Unsetenv("MAILJET_API_SECRET")
		_ = os.Unsetenv("F2K_API_KEY")
	}()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "test@kindle.com", cfg.KindleEmail)
	assert.Equal(t, "sender@example.com", cfg.SenderEmail)
	assert.Equal(t, "api-key", cfg.MailjetAPIKey)
	assert.Equal(t, "api-secret", cfg.MailjetAPISecret)
	assert.Equal(t, "api-key-secret", cfg.APIKeySecret)
}
