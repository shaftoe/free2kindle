// Package auth provides pluggable authentication backends for the free2kindle application.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	userIDKey    contextKey = "user_id"
	authErrorKey contextKey = "auth_error"
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

// GetAuthError retrieves the authentication error from the context, if any.
func GetAuthError(ctx context.Context) error {
	authError, found := ctx.Value(authErrorKey).(string)
	if found {
		return errors.New(authError)
	}
	return nil
}

func sharedAPIKeyMiddleware(apiKeySecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get(authHeader)
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				handleAuthError(r.Context(), next, w, r, "Missing or malformed auth header")
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			if token != apiKeySecret {
				handleAuthError(r.Context(), next, w, r, "Invalid API key")
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

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get(authHeader)
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				ctx := addUserIDToContext(r.Context(), anonymousUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			claims, validateErr := jwtValidator.ValidateToken(r.Context(), token)
			if validateErr != nil {
				handleAuthError(r.Context(), next, w, r, "Invalid JWT token: "+validateErr.Error())
				return
			}

			validatedClaims, ok := claims.(*validator.ValidatedClaims)
			if !ok {
				handleAuthError(r.Context(), next, w, r, "Failed to parse JWT claims")
				return
			}

			ctx := addUserIDToContext(r.Context(), validatedClaims.RegisteredClaims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// EnsureAutheticatedMiddleware ensures that the request is authenticated before
// proceeding to the next handler.
func EnsureAutheticatedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := GetAuthError(r.Context()); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
			return
		}

		accountID := GetAccountID(r.Context())
		if accountID == "" || accountID == anonymousUserID {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "Unauthorized (missing account ID)"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func addUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func handleAuthError(ctx context.Context, next http.Handler, w http.ResponseWriter, r *http.Request, msg string) {
	ctx = context.WithValue(ctx, authErrorKey, msg)

	next.ServeHTTP(w, r.WithContext(ctx))
}
