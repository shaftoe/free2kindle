// Package mailjet provides a Mailjet implementation of the email Sender interface.
package mailjet

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	mailjetLib "github.com/mailjet/mailjet-apiv3-go/v4"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/email"
)

// Config holds the Mailjet API configuration.
type Config struct {
	APIKey      string
	APISecret   string
	SenderEmail string
}

// Sender implements the email.Sender interface using Mailjet.
type Sender struct {
	config *Config
	client *mailjetLib.Client
}

// NewSender creates a new Mailjet sender instance.
func NewSender(config *Config) *Sender {
	mailjetClient := mailjetLib.NewMailjetClient(config.APIKey, config.APISecret)
	return &Sender{
		config: config,
		client: mailjetClient,
	}
}

// SendEmail sends an email with the EPUB attachment via Mailjet.
func (s *Sender) SendEmail(_ context.Context, req *email.EmailRequest) (*email.SendEmailResponse, error) {
	if err := s.validateConfig(); err != nil {
		return nil, fmt.Errorf("invalid sender config: %w", err)
	}

	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	filename := email.GenerateFilename(req.Article)
	subject := email.GenerateSubject(req.Article.Title, req.Subject)

	base64Content := base64.StdEncoding.EncodeToString(req.EPUBData)

	messagesInfo := []mailjetLib.InfoMessagesV31{
		{
			From: &mailjetLib.RecipientV31{
				Email: s.config.SenderEmail,
			},
			To: &mailjetLib.RecipientsV31{
				mailjetLib.RecipientV31{
					Email: req.KindleEmail,
				},
			},
			Subject:  subject,
			TextPart: "EPUB document attached.",
			Attachments: &mailjetLib.AttachmentsV31{
				mailjetLib.AttachmentV31{
					ContentType:   "application/epub+zip",
					Filename:      filename,
					Base64Content: base64Content,
				},
			},
		},
	}

	messages := mailjetLib.MessagesV31{Info: messagesInfo}

	resp, err := s.client.SendMailV31(&messages)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	if len(resp.ResultsV31) == 0 {
		return nil, errors.New("no messages in response")
	}

	if resp.ResultsV31[0].Status != "success" {
		return nil, fmt.Errorf("email send failed with status: %s", resp.ResultsV31[0].Status)
	}

	result := &email.SendEmailResponse{
		Status:  "success",
		Message: "Email sent successfully",
	}

	if len(resp.ResultsV31[0].To) > 0 {
		result.EmailUUID = resp.ResultsV31[0].To[0].MessageUUID
	}

	return result, nil
}

func (s *Sender) validateConfig() error {
	if s.config.APIKey == "" {
		return errors.New("API key is required")
	}
	if s.config.APISecret == "" {
		return errors.New("API secret is required")
	}
	if s.config.SenderEmail == "" {
		return errors.New("sender email is required")
	}
	return nil
}

func (s *Sender) validateRequest(req *email.EmailRequest) error {
	if req.KindleEmail == "" {
		return errors.New("kindle email is required")
	}
	if req.EPUBData == nil {
		return errors.New("EPUB data is required")
	}
	if len(req.EPUBData) == 0 {
		return errors.New("EPUB data is empty")
	}
	if req.Article == nil {
		return errors.New("article is required")
	}
	return nil
}
