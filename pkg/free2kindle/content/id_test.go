package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleIDFromURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid URL",
			url:     "https://example.com/article/123",
			wantErr: false,
		},
		{
			name:    "valid URL with query params",
			url:     "https://example.com/article/123?source=twitter&utm=test",
			wantErr: false,
		},
		{
			name:    "valid URL with fragment",
			url:     "https://example.com/article/123#section-1",
			wantErr: false,
		},
		{
			name:    "valid URL with both query and fragment",
			url:     "https://example.com/article/123?ref=news#intro",
			wantErr: false,
		},
		{
			name:    "valid HTTP URL",
			url:     "http://example.com/article/456",
			wantErr: false,
		},
		{
			name:    "URL with trailing slash",
			url:     "https://example.com/article/123/",
			wantErr: false,
		},
		{
			name:        "invalid URL",
			url:         "not-a-url",
			wantErr:     true,
			errContains: "must have scheme and host",
		},
		{
			name:        "URL without scheme",
			url:         "example.com/article",
			wantErr:     true,
			errContains: "must have scheme and host",
		},
		{
			name:        "empty URL",
			url:         "",
			wantErr:     true,
			errContains: "must have scheme and host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ArticleIDFromURL(tt.url)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Empty(t, id)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, id)
			assert.Len(t, id, 36)
		})
	}
}

func TestArticleIDFromURL_Deterministic(t *testing.T) {
	url1 := "https://example.com/article/123?source=twitter"
	url2 := "https://example.com/article/123?utm_source=newsletter#intro"
	url3 := "https://example.com/article/123/"

	id1, err1 := ArticleIDFromURL(url1)
	id2, err2 := ArticleIDFromURL(url2)
	id3, err3 := ArticleIDFromURL(url3)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	assert.Equal(t, id1, id2, "IDs should be same for same base URL")
	assert.Equal(t, id1, id3, "IDs should be same for same base URL")
}

func TestArticleIDFromURL_DifferentURLs(t *testing.T) {
	urls := []string{
		"https://example.com/article/1",
		"https://example.com/article/2",
		"https://other.com/article/1",
		"https://example.org/article/1",
	}

	ids := make(map[string]bool)
	for _, u := range urls {
		id, err := ArticleIDFromURL(u)
		assert.NoError(t, err)
		assert.False(t, ids[id], "ID should be unique for each URL")
		ids[id] = true
	}

	assert.Equal(t, len(urls), len(ids))
}

func TestArticleIDFromURL_HttpVsHttps(t *testing.T) {
	idHTTP, err1 := ArticleIDFromURL("http://example.com/article")
	idHTTPS, err2 := ArticleIDFromURL("https://example.com/article")

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, idHTTP, idHTTPS, "HTTP and HTTPS should produce different IDs")
}
