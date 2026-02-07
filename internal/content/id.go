package content

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

// ArticleIDFromURL generates a deterministic UUID v5 for an article from its URL.
// Strips query parameters and fragments before hashing to ensure
// the same base URL always produces the same ID.
//
// Uses UUID v5 with the URL namespace as defined in RFC 4122.
func ArticleIDFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("URL must be valid: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("URL must have scheme and host")
	}

	path := strings.TrimSuffix(parsedURL.Path, "/")
	if path == "" {
		path = "/"
	}

	cleanURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, path)

	id := uuid.NewSHA1(uuid.NameSpaceURL, []byte(cleanURL))

	return id.String(), nil
}
