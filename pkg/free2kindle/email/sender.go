// Package email provides email sending functionality for Kindle devices.
package email

import (
	"context"
	"regexp"
	"strings"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
)

// SendEmailResponse contains the response from sending an email.
type SendEmailResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	EmailUUID string `json:"email_uuid,omitempty"`
}

// Sender defines the interface for sending emails.
type Sender interface {
	SendEmail(ctx context.Context, req *Request) (*SendEmailResponse, error)
}

// Request contains the data required to send an email.
type Request struct {
	Article     *content.Article
	EPUBData    []byte
	KindleEmail string
	Subject     string
}

// GenerateFilename creates a sanitized filename from the article title.
func GenerateFilename(article *content.Article) string {
	if article.Title != "" {
		return sanitizeFilename(article.Title) + ".epub"
	}
	return "article.epub"
}

// GenerateSubject creates an email subject from the article title or custom subject.
func GenerateSubject(articleTitle, customSubject string) string {
	if customSubject != "" {
		return sanitizeSubject(customSubject)
	}
	if articleTitle != "" {
		return sanitizeSubject(articleTitle)
	}
	return "Document"
}

func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[^\w\s-]`)
	sanitized := re.ReplaceAllString(name, "")
	sanitized = strings.TrimSpace(sanitized)
	if sanitized == "" {
		return "article"
	}
	return sanitized
}

func sanitizeSubject(subject string) string {
	if subject == "" {
		return "Document"
	}
	subject = strings.TrimSpace(subject)
	if len(subject) > 100 {
		subject = subject[:100]
	}
	return subject
}
