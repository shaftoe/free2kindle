// Package content provides article extraction functionality from web pages.
package content

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-trafilatura"
)

const (
	wordsPerMinute = 250
)

// Article represents the extracted content from a web page.
type Article struct {
	ID                 string
	Title              string
	Author             string
	Content            string
	Excerpt            string
	URL                string
	ImageURL           string
	PublishedAt        time.Time
	HTML               string
	ExtractedAt        time.Time
	WordCount          int
	ReadingTimeMinutes int
	SourceDomain       string
	SiteName           string
	ContentType        string
	Language           string
}

// Extractor handles the extraction of article content from URLs and HTML.
type Extractor struct {
	client *http.Client
}

// NewExtractor creates a new Extractor instance.
func NewExtractor() *Extractor {
	return &Extractor{
		client: &http.Client{},
	}
}

// ExtractFromURL fetches and extracts article content from the given URL.
func (e *Extractor) ExtractFromURL(ctx context.Context, urlStr string) (*Article, error) {
	id, err := ArticleIDFromURL(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article ID: %w", err)
	}

	parsedURL, body, err := e.fetchURL(ctx, urlStr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := body.Close(); closeErr != nil {
			log.Printf("warning: failed to close response body: %v", closeErr)
		}
	}()

	opts := trafilatura.Options{
		OriginalURL: parsedURL,
	}

	result, err := trafilatura.Extract(body, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article content: %w", err)
	}

	if result.ContentNode == nil {
		return nil, errors.New("no content extracted")
	}

	return e.buildArticle(result, urlStr, id), nil
}

func (e *Extractor) fetchURL(ctx context.Context, urlStr string) (*url.URL, io.ReadCloser, error) {
	if err := validateURL(urlStr); err != nil {
		return nil, nil, fmt.Errorf("invalid URL: %w", err)
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch URL: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("warning: failed to close response body: %v", closeErr)
		}
		return nil, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("warning: failed to close response body: %v", closeErr)
		}
		return nil, nil, fmt.Errorf("expected HTML content, got: %s", contentType)
	}

	return parsedURL, resp.Body, nil
}

func (e *Extractor) buildArticle(result *trafilatura.ExtractResult, urlStr, id string) *Article {
	contentHTML := dom.InnerHTML(result.ContentNode)
	plainText := stripHTML(contentHTML)
	wordCount := countWords(plainText)

	return &Article{
		ID:                 id,
		Title:              result.Metadata.Title,
		Author:             result.Metadata.Author,
		Content:            contentHTML,
		Excerpt:            result.Metadata.Description,
		ImageURL:           result.Metadata.Image,
		PublishedAt:        result.Metadata.Date,
		URL:                urlStr,
		ExtractedAt:        time.Now(),
		WordCount:          wordCount,
		ReadingTimeMinutes: (wordCount + wordsPerMinute - 1) / wordsPerMinute,
		SourceDomain:       result.Metadata.Hostname,
		SiteName:           result.Metadata.Sitename,
		ContentType:        result.Metadata.PageType,
		Language:           result.Metadata.Language,
	}
}

func stripHTML(html string) string {
	re := strings.NewReplacer(
		"<p>", " ",
		"</p>", " ",
		"<div>", " ",
		"</div>", " ",
		"<br>", " ",
		"<br/>", " ",
		"<br />", " ",
	)

	result := re.Replace(html)

	var stripped strings.Builder
	inTag := false
	for _, r := range result {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				stripped.WriteRune(r)
			}
		}
	}

	return stripped.String()
}

func countWords(text string) int {
	fields := strings.Fields(text)
	return len(fields)
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
