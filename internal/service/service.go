// Package service provides the main orchestration logic for processing articles.
package service

import (
	"context"
	"fmt"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/email"
	"github.com/shaftoe/free2kindle/internal/epub"
	"github.com/shaftoe/free2kindle/internal/model"
)

// Deps holds the external dependencies required by the service.
type Deps struct {
	extractor *content.Extractor
	generator *epub.Generator
	sender    email.Sender
}

// Options holds runtime options for the service execution.
type Options struct {
	sendEmail    bool
	generateEPUB bool
	subject      string
	outputPath   string
}

// NewDeps creates a new Deps struct with the given components.
func NewDeps(extractor *content.Extractor, generator *epub.Generator, sender email.Sender) *Deps {
	return &Deps{
		extractor: extractor,
		generator: generator,
		sender:    sender,
	}
}

// NewOptions creates a new Options struct with the given parameters.
func NewOptions(sendEmail, generateEPUB bool, subject, outputPath string) *Options {
	return &Options{
		sendEmail:    sendEmail,
		generateEPUB: generateEPUB,
		subject:      subject,
		outputPath:   outputPath,
	}
}

// Result contains the output from processing an article.
type Result struct {
	Article  *model.Article
	EPUBData []byte
	URL      string
}

// Run processes a URL to extract content, generate EPUB, and optionally send email.
func Run(ctx context.Context, d *Deps, cfg *config.Config, opts *Options, url string) (*Result, error) {
	article, err := d.extractor.ExtractFromURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article: %w", err)
	}

	if article.Title == "" {
		article.Title = "Untitled"
	}

	var epubData []byte
	if opts.sendEmail || opts.generateEPUB {
		epubData, err = d.generator.Generate(article)
		if err != nil {
			return nil, fmt.Errorf("failed to generate EPUB: %w", err)
		}
	}

	if opts.sendEmail {
		emailReq := &email.Request{
			Article:   article,
			EPUBData:  epubData,
			DestEmail: cfg.DestEmail,
			Subject:   email.GenerateSubject(article.Title, opts.subject),
		}

		if _, sendErr := d.sender.SendEmail(ctx, emailReq); sendErr != nil {
			return nil, fmt.Errorf("failed to send email: %w", sendErr)
		}
	}

	if opts.generateEPUB && opts.outputPath != "" {
		if writeErr := d.generator.GenerateAndWrite(article, opts.outputPath); writeErr != nil {
			return nil, fmt.Errorf("failed to write EPUB document: %w", writeErr)
		}
	}

	return &Result{
		Article:  article,
		EPUBData: epubData,
	}, nil
}
