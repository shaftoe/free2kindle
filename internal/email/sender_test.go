package email

import (
	"testing"

	"github.com/shaftoe/savetoink/internal/model"
)

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		name     string
		article  *model.Article
		expected string
	}{
		{
			name: "article with title",
			article: &model.Article{
				Title: "Test Article",
			},
			expected: "Test Article.epub",
		},
		{
			name: "article with special characters in title",
			article: &model.Article{
				Title: "Test Article: What's New?",
			},
			expected: "Test Article Whats New.epub",
		},
		{
			name:     "article without title",
			article:  &model.Article{},
			expected: "article.epub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateFilename(tt.article)
			if got != tt.expected {
				t.Errorf("GenerateFilename() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple filename",
			input:    "test-file",
			expected: "test-file",
		},
		{
			name:     "filename with special chars",
			input:    "test: file? what!",
			expected: "test file what",
		},
		{
			name:     "empty filename",
			input:    "",
			expected: "article",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "article",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeFilename() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGenerateSubject(t *testing.T) {
	tests := []struct {
		name          string
		articleTitle  string
		customSubject string
		expected      string
	}{
		{
			name:          "custom subject",
			articleTitle:  "Article Title",
			customSubject: "Custom Subject",
			expected:      "Custom Subject",
		},
		{
			name:          "article title subject",
			articleTitle:  "Article Title",
			customSubject: "",
			expected:      "Article Title",
		},
		{
			name:          "empty",
			articleTitle:  "",
			customSubject: "",
			expected:      "Document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSubject(tt.articleTitle, tt.customSubject)
			if got != tt.expected {
				t.Errorf("GenerateSubject() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSanitizeSubject(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple subject",
			input:    "Test Subject",
			expected: "Test Subject",
		},
		{
			name:     "subject with leading/trailing spaces",
			input:    "  Test Subject  ",
			expected: "Test Subject",
		},
		{
			name: "long subject",
			input: "This is a very long subject that definitely exceeds the " +
				"100 character limit and should be properly truncated",
			expected: "This is a very long subject that definitely exceeds " +
				"the 100 character limit and should be properly t",
		},
		{
			name:     "empty subject",
			input:    "",
			expected: "Document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeSubject(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeSubject() = %v, want %v", got, tt.expected)
			}
		})
	}
}
