// Package auth provides pluggable authentication backends for the free2kindle application.
package auth

import (
	"encoding/json"
	"net/http"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/constant"
)

const (
	apiKeyHeader = "X-API-Key" //nolint:gosec // G101: This is a header name constant, not a hardcoded credential
)

type errorResponse struct {
	Message string `json:"message"`
}

// NewMiddleware returns authentication middleware based on the configured auth backend.
func NewMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	switch cfg.AuthBackend {
	case constant.AuthBackendSharedAPIKey:
		return sharedAPIKeyMiddleware(cfg.APIKeySecret)
	default:
		return sharedAPIKeyMiddleware(cfg.APIKeySecret)
	}
}

func sharedAPIKeyMiddleware(apiKeySecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(apiKeyHeader)
			if apiKey == "" || apiKey != apiKeySecret {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(errorResponse{Message: "Invalid API key"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
