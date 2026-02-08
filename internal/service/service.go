// Package service provides the main orchestration logic for processing articles.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/email"
	"github.com/shaftoe/free2kindle/internal/epub"
	"github.com/shaftoe/free2kindle/internal/model"
)

// Interface defines the contract for service operations.
type Interface interface {
	Process(ctx context.Context, url string) (*ProcessResult, error)
	Send(ctx context.Context, result *ProcessResult, subject string) (*email.SendEmailResponse, error)
	WriteToFile(result *ProcessResult, outputPath string) error
}

// Deps holds the external dependencies required by the service.
type Deps struct {
	Extractor *content.Extractor
	Generator *epub.Generator
	Sender    email.Sender
}

// NewDeps creates a new Deps struct with the given components.
func NewDeps(extractor *content.Extractor, generator *epub.Generator, sender email.Sender) *Deps {
	return &Deps{
		Extractor: extractor,
		Generator: generator,
		Sender:    sender,
	}
}

// Service holds stateless dependencies and provides methods to process articles.
type Service struct {
	extractor *content.Extractor
	generator *epub.Generator
	sender    email.Sender
	cfg       *config.Config
}

// New creates a new Service instance with the given dependencies.
func New(d *Deps, cfg *config.Config) *Service {
	return &Service{
		extractor: d.Extractor,
		generator: d.Generator,
		sender:    d.Sender,
		cfg:       cfg,
	}
}

// ProcessResult holds the result of processing an article.
type ProcessResult struct {
	article  *model.Article
	epubData []byte
	url      string
}

// Article returns the extracted article.
func (r *ProcessResult) Article() *model.Article {
	return r.article
}

// EPUBData returns the generated EPUB data.
func (r *ProcessResult) EPUBData() []byte {
	return r.epubData
}

// URL returns the URL that was processed.
func (r *ProcessResult) URL() string {
	return r.url
}

// NewProcessResult creates a new ProcessResult for testing purposes.
// This is primarily used in tests to create mock results.
func NewProcessResult(article *model.Article, epubData []byte, url string) *ProcessResult {
	return &ProcessResult{
		article:  article,
		epubData: epubData,
		url:      url,
	}
}

// Process extracts content from a URL and generates EPUB data.
// Can be called multiple times to re-fetch fresh content.
func (s *Service) Process(ctx context.Context, url string) (*ProcessResult, error) {
	article, err := s.extractor.ExtractFromURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article: %w", err)
	}

	if article.Title == "" {
		article.Title = "Untitled"
	}

	epubData, err := s.generator.Generate(article)
	if err != nil {
		return nil, fmt.Errorf("failed to generate EPUB: %w", err)
	}

	return &ProcessResult{
		article:  article,
		epubData: epubData,
		url:      url,
	}, nil
}

// Send sends an email with the processed article and EPUB.
// Returns an error if the result is nil or if sending fails.
// Can be called multiple times with the same result.
func (s *Service) Send(
	ctx context.Context,
	result *ProcessResult,
	subject string,
) (*email.SendEmailResponse, error) {
	if result == nil {
		return nil, errors.New("result is nil, must call Process first")
	}

	if result.article == nil {
		return nil, errors.New("article is nil, must call Process first")
	}

	if s.sender == nil {
		return nil, errors.New("email sender is not configured")
	}

	emailReq := &email.Request{
		Article:   result.article,
		EPUBData:  result.epubData,
		DestEmail: s.cfg.DestEmail,
		Subject:   email.GenerateSubject(result.article.Title, subject),
	}

	resp, err := s.sender.SendEmail(ctx, emailReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return resp, nil
}

// WriteToFile writes the EPUB data to a file.
// Returns an error if the result is nil or if writing fails.
func (s *Service) WriteToFile(result *ProcessResult, outputPath string) error {
	if result == nil {
		return errors.New("result is nil, must call Process first")
	}

	if result.article == nil {
		return errors.New("article is nil, must call Process first")
	}

	if outputPath == "" {
		return errors.New("output path is empty")
	}

	err := s.generator.GenerateAndWrite(result.article, outputPath)
	if err != nil {
		return fmt.Errorf("failed to write EPUB document: %w", err)
	}

	return nil
}
