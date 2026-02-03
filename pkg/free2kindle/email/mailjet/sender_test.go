package mailjet

import (
	"context"
	"testing"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email"
)

func TestNewSender(t *testing.T) {
	config := &Config{
		APIKey:      "test-key",
		APISecret:   "test-secret",
		SenderEmail: "test@example.com",
	}

	sender := NewSender(config)
	if sender == nil { //nolint:staticcheck
		t.Fatal("NewSender returned nil")
	}

	if sender.config != config { //nolint:staticcheck
		t.Error("Sender config not set correctly")
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				APIKey:      "key",
				APISecret:   "secret",
				SenderEmail: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: &Config{
				APISecret:   "secret",
				SenderEmail: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing api secret",
			config: &Config{
				APIKey:      "key",
				SenderEmail: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing sender email",
			config: &Config{
				APIKey:    "key",
				APISecret: "secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := NewSender(tt.config)
			err := sender.validateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *email.EmailRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &email.EmailRequest{
				Article: &content.Article{
					Title: "Test Article",
				},
				EPUBData:    []byte("test epub data"),
				KindleEmail: "kindle@kindle.com",
			},
			wantErr: false,
		},
		{
			name: "missing kindle email",
			req: &email.EmailRequest{
				Article:     &content.Article{Title: "Test"},
				EPUBData:    []byte("data"),
				KindleEmail: "",
			},
			wantErr: true,
		},
		{
			name: "missing epub data",
			req: &email.EmailRequest{
				Article:     &content.Article{Title: "Test"},
				EPUBData:    nil,
				KindleEmail: "kindle@kindle.com",
			},
			wantErr: true,
		},
		{
			name: "missing article",
			req: &email.EmailRequest{
				Article:     nil,
				EPUBData:    []byte("data"),
				KindleEmail: "kindle@kindle.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := NewSender(&Config{
				APIKey:      "key",
				APISecret:   "secret",
				SenderEmail: "test@example.com",
			})
			err := sender.validateRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendEmailValidation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		config  *Config
		req     *email.EmailRequest
		wantErr bool
	}{
		{
			name: "missing config",
			config: &Config{
				APIKey:      "",
				APISecret:   "secret",
				SenderEmail: "test@example.com",
			},
			req: &email.EmailRequest{
				Article:     &content.Article{Title: "Test"},
				EPUBData:    []byte("data"),
				KindleEmail: "kindle@kindle.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := NewSender(tt.config)
			err := sender.SendEmail(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
