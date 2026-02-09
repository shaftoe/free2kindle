// Package auth provides pluggable authentication backends for the free2kindle application.
package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/constant"
)

const (
	apiKeyHeader   = "X-API-Key" //nolint:gosec // G101: This is a header name constant, not a hardcoded credential
	adminAccountID = "admin"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
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

// GetAccountID retrieves the authenticated account ID from the context.
func GetAccountID(ctx context.Context) string {
	accountID, _ := ctx.Value(userIDKey).(string)
	return accountID
}

func sharedAPIKeyMiddleware(apiKeySecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(apiKeyHeader)
			if apiKey == "" || apiKey != apiKeySecret {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, adminAccountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func EnsureAutheticatedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID := GetAccountID(r.Context())
		if accountID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Message: "Unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
