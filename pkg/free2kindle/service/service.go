// Package service provides the main orchestration logic for processing articles.
package service

import (
	"context"
	"fmt"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
)

// Config holds the configuration for running the service.
type Config struct {
	Extractor    *content.Extractor
	Generator    *epub.Generator
	Sender       email.Sender
	SendEmail    bool
	GenerateEPUB bool
	KindleEmail  string
	SenderEmail  string
	Subject      string
	OutputPath   string
}

// Result contains the output from processing an article.
type Result struct {
	Article  *content.Article
	EPUBData []byte
	Title    string
	URL      string
}

// Run processes a URL to extract content, generate EPUB, and optionally send email.
func Run(ctx context.Context, cfg *Config, url string) (*Result, error) {
	article, err := cfg.Extractor.ExtractFromURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article: %w", err)
	}

	if article.Title == "" {
		article.Title = "Untitled"
	}

	var epubData []byte
	if cfg.SendEmail || cfg.GenerateEPUB {
		epubData, err = cfg.Generator.Generate(article)
		if err != nil {
			return nil, fmt.Errorf("failed to generate EPUB: %w", err)
		}
	}

	if cfg.SendEmail {
		emailReq := &email.EmailRequest{
			Article:     article,
			EPUBData:    epubData,
			KindleEmail: cfg.KindleEmail,
			Subject:     email.GenerateSubject(article.Title, cfg.Subject),
		}

		if _, err := cfg.Sender.SendEmail(ctx, emailReq); err != nil {
			return nil, fmt.Errorf("failed to send email: %w", err)
		}
	}

	if cfg.GenerateEPUB && cfg.OutputPath != "" {
		if err := cfg.Generator.GenerateAndWrite(article, cfg.OutputPath); err != nil {
			return nil, fmt.Errorf("failed to write EPUB: %w", err)
		}
	}

	return &Result{
		Article:  article,
		EPUBData: epubData,
		Title:    article.Title,
		URL:      url,
	}, nil
}
