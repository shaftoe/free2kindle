// Package email provides email sending functionality.
package email

import (
	"context"
	"regexp"
	"strings"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
)

const (
	// DefaultSubject is the default email subject.
	DefaultSubject = "Document"
	// MaxSubjectLength is the maximum length for email subjects.
	MaxSubjectLength = 100
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
	// Article is the article to be sent.
	Article *content.Article

	// EPUBData is the EPUB data to be sent as attachment.
	EPUBData []byte

	// Subject is the email subject.
	Subject string

	// DestEmail is the email address of the recipient, typically a
	// Kindle Personal Document Service address like "abcd@kindle.com".
	DestEmail string
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
	return DefaultSubject
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
		return DefaultSubject
	}
	subject = strings.TrimSpace(subject)
	if len(subject) > MaxSubjectLength {
		subject = subject[:MaxSubjectLength]
	}
	return subject
}
