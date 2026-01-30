package mailjet

import (
	"context"
	"encoding/base64"
	"fmt"

	mailjetLib "github.com/mailjet/mailjet-apiv3-go/v4"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/email"
)

type Config struct {
	APIKey      string
	APISecret   string
	SenderEmail string
}

type Sender struct {
	config *Config
	client *mailjetLib.Client
}

func NewSender(config *Config) *Sender {
	mailjetClient := mailjetLib.NewMailjetClient(config.APIKey, config.APISecret)
	return &Sender{
		config: config,
		client: mailjetClient,
	}
}

func (s *Sender) SendEmail(ctx context.Context, req *email.EmailRequest) error {
	if err := s.validateConfig(); err != nil {
		return fmt.Errorf("invalid sender config: %w", err)
	}

	if err := s.validateRequest(req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
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
		return fmt.Errorf("failed to send email: %w", err)
	}

	if len(resp.ResultsV31) == 0 {
		return fmt.Errorf("no messages in response")
	}

	if resp.ResultsV31[0].Status != "success" {
		return fmt.Errorf("email send failed with status: %s", resp.ResultsV31[0].Status)
	}

	if len(resp.ResultsV31[0].To) > 0 {
		fmt.Printf("Email sent successfully. Message ID: %d, UUID: %s\n",
			resp.ResultsV31[0].To[0].MessageID, resp.ResultsV31[0].To[0].MessageUUID)
	} else {
		fmt.Println("Email sent successfully")
	}

	return nil
}

func (s *Sender) validateConfig() error {
	if s.config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	if s.config.APISecret == "" {
		return fmt.Errorf("API secret is required")
	}
	if s.config.SenderEmail == "" {
		return fmt.Errorf("sender email is required")
	}
	return nil
}

func (s *Sender) validateRequest(req *email.EmailRequest) error {
	if req.KindleEmail == "" {
		return fmt.Errorf("kindle email is required")
	}
	if req.EPUBData == nil {
		return fmt.Errorf("EPUB data is required")
	}
	if len(req.EPUBData) == 0 {
		return fmt.Errorf("EPUB data is empty")
	}
	if req.Article == nil {
		return fmt.Errorf("article is required")
	}
	return nil
}
