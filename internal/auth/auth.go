// Package auth provides pluggable authentication backends for the free2kindle application.
package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/jwks"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/constant"
)

const (
	authHeader       = "Authorization"
	adminAccountID   = "admin"
	anonymousUserID  = "-"
	allowedClockSkew = 30 * time.Second
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
)

type errorResponse struct {
	Message string `json:"message"`
}

// NewUserIDMiddleware returns authentication middleware based on the configured auth backend.
// Ensure the userID is set in the context (`anonymousUserID` string for anonymous users).
func NewUserIDMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	switch cfg.AuthBackend {
	case constant.AuthBackendSharedAPIKey:
		return sharedAPIKeyMiddleware(cfg.APIKeySecret)
	case constant.AuthBackendAuth0:
		return auth0Middleware(cfg.Auth0Domain, cfg.Auth0Audience)
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
			auth := r.Header.Get(authHeader)
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				ctx := addUserIDToContext(r.Context(), anonymousUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			if token != apiKeySecret {
				ctx := addUserIDToContext(r.Context(), anonymousUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			ctx := addUserIDToContext(r.Context(), adminAccountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func auth0Middleware(domain, audience string) func(http.Handler) http.Handler {
	issuerURL, err := url.Parse("https://" + domain + "/")
	if err != nil {
		panic("failed to parse issuer URL: " + err.Error())
	}

	provider, err := jwks.NewCachingProvider(
		jwks.WithIssuerURL(issuerURL),
	)
	if err != nil {
		panic("failed to create JWKS provider: " + err.Error())
	}

	jwtValidator, err := validator.New(
		validator.WithKeyFunc(provider.KeyFunc),
		validator.WithAlgorithm(validator.RS256),
		validator.WithIssuer(issuerURL.String()),
		validator.WithAudience(audience),
		validator.WithAllowedClockSkew(allowedClockSkew),
	)
	if err != nil {
		panic("failed to create JWT validator: " + err.Error())
	}

	middleware, err := jwtmiddleware.New(
		jwtmiddleware.WithValidator(jwtValidator),
		jwtmiddleware.WithCredentialsOptional(true),
	)
	if err != nil {
		panic("failed to create JWT middleware: " + err.Error())
	}

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, claimsErr := jwtmiddleware.GetClaims[*validator.ValidatedClaims](r.Context())
			if claimsErr != nil {
				ctx := addUserIDToContext(r.Context(), anonymousUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := addUserIDToContext(r.Context(), claims.RegisteredClaims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		}))
	}
}

// EnsureAutheticatedMiddleware ensures that the request is authenticated before
// proceeding to the next handler.
func EnsureAutheticatedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID := GetAccountID(r.Context())
		if accountID == "" || accountID == anonymousUserID {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "Unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func addUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
