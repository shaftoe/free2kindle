package content

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-trafilatura"
)

type Article struct {
	Title       string
	Author      string
	Content     string
	Excerpt     string
	URL         string
	ImageURL    string
	PublishedAt time.Time
	HTML        string
}

type Extractor struct {
	client *http.Client
}

func NewExtractor() *Extractor {
	return &Extractor{
		client: &http.Client{},
	}
}

func (e *Extractor) ExtractFromURL(ctx context.Context, urlStr string) (*Article, error) {
	if err := validateURL(urlStr); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("warning: failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("expected HTML content, got: %s", contentType)
	}

	opts := trafilatura.Options{
		OriginalURL: parsedURL,
	}

	result, err := trafilatura.Extract(resp.Body, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article content: %w", err)
	}

	if result.ContentNode == nil {
		return nil, errors.New("no content extracted")
	}

	contentHTML := dom.InnerHTML(result.ContentNode)

	return &Article{
		Title:       result.Metadata.Title,
		Author:      result.Metadata.Author,
		Content:     contentHTML,
		Excerpt:     result.Metadata.Description,
		ImageURL:    result.Metadata.Image,
		PublishedAt: result.Metadata.Date,
		URL:         urlStr,
		HTML:        "",
	}, nil
}

func (e *Extractor) ExtractFromHTML(ctx context.Context, urlStr, html string) (*Article, error) {
	if err := validateURL(urlStr); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	opts := trafilatura.Options{
		OriginalURL: parsedURL,
	}

	result, err := trafilatura.Extract(strings.NewReader(html), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article content: %w", err)
	}

	if result.ContentNode == nil {
		return nil, errors.New("no content extracted")
	}

	contentHTML := dom.InnerHTML(result.ContentNode)

	return &Article{
		Title:       result.Metadata.Title,
		Author:      result.Metadata.Author,
		Content:     contentHTML,
		Excerpt:     result.Metadata.Description,
		ImageURL:    result.Metadata.Image,
		PublishedAt: result.Metadata.Date,
		URL:         urlStr,
		HTML:        html,
	}, nil
}

func validateURL(urlStr string) error {
	if urlStr == "" {
		return errors.New("URL cannot be empty")
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("URL must use http or https scheme")
	}

	if u.Host == "" {
		return errors.New("URL must have a host")
	}

	return nil
}
