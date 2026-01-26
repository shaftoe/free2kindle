package email

import (
	"context"
	"regexp"
	"strings"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
)

type Sender interface {
	SendEmail(ctx context.Context, req *EmailRequest) error
}

type EmailRequest struct {
	Article     *content.Article
	EPUBData    []byte
	KindleEmail string
	Subject     string
}

func GenerateFilename(article *content.Article) string {
	if article.Title != "" {
		return sanitizeFilename(article.Title) + ".epub"
	}
	return "article.epub"
}

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
