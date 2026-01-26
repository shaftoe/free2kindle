package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email/mailjet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func sendToKindle(ctx context.Context, article *content.Article) error {
	kindleEmail := viper.GetString("kindle-email")
	senderEmail := viper.GetString("sender-email")
	apiKey := viper.GetString("api-key")
	apiSecret := viper.GetString("api-secret")

	config := &mailjet.Config{
		APIKey:      apiKey,
		APISecret:   apiSecret,
		SenderEmail: senderEmail,
	}

	sender := mailjet.NewSender(config)

	fmt.Printf("Generating EPUB for email...\n")
	start := time.Now()
	epubData, err := generator.Generate(article)
	if err != nil {
		return fmt.Errorf("failed to generate EPUB: %w", err)
	}
	fmt.Printf("Generated in %v\n", time.Since(start))

	emailReq := &email.EmailRequest{
		Article:     article,
		EPUBData:    epubData,
		KindleEmail: kindleEmail,
		Subject:     emailSubject,
	}

	fmt.Printf("Sending to Kindle: %s -> %s\n", senderEmail, kindleEmail)
	start = time.Now()
	if err := sender.SendEmail(ctx, emailReq); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	fmt.Printf("Sent in %v\n", time.Since(start))

	fmt.Printf("\n✓ Article sent to Kindle\n")

	return nil
}

func validateEmailConfig(cmd *cobra.Command, args []string) error {
	if !sendEmail {
		return nil
	}

	kindleEmail := viper.GetString("kindle-email")
	senderEmail := viper.GetString("sender-email")
	apiKey := viper.GetString("api-key")
	apiSecret := viper.GetString("api-secret")

	var missing []string
	if kindleEmail == "" {
		missing = append(missing, "--kindle-email or F2K_KINDLE_EMAIL")
	}
	if senderEmail == "" {
		missing = append(missing, "--sender-email or F2K_SENDER_EMAIL")
	}
	if apiKey == "" {
		missing = append(missing, "--api-key or MAILJET_API_KEY")
	}
	if apiSecret == "" {
		missing = append(missing, "--api-secret or MAILJET_API_SECRET")
	}

	if len(missing) > 0 {
		return fmt.Errorf("when using --send flag, the following are required: %s",
			stringJoin(missing, ", "))
	}

	return nil
}
