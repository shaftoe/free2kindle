package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/constant"
)

func TestNewMiddleware_SharedAPIKey(t *testing.T) {
	cfg := &config.Config{
		AuthBackend:  constant.AuthBackendSharedAPIKey,
		APIKeySecret: "valid-key",
	}

	middleware := NewUserIDMiddleware(cfg)

	if middleware == nil {
		t.Fatal("expected middleware to not be nil")
	}
}

func TestNewMiddleware_Auth0(t *testing.T) {
	cfg := &config.Config{
		AuthBackend:   constant.AuthBackendAuth0,
		Auth0Domain:   "example.auth0.com",
		Auth0Audience: "test-audience",
	}

	middleware := NewUserIDMiddleware(cfg)

	if middleware == nil {
		t.Fatal("expected middleware to not be nil")
	}
}

func TestSharedAPIKeyMiddleware_ValidAPIKey(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer valid-key")
	w := httptest.NewRecorder()

	sharedAPIKeyMiddleware("valid-key")(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSharedAPIKeyMiddleware_MissingAPIKey(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	w := httptest.NewRecorder()

	middlewareChain := EnsureAutheticatedMiddleware(sharedAPIKeyMiddleware("valid-key")(next))
	middlewareChain.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestSharedAPIKeyMiddleware_WrongAPIKey(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer wrong-key")
	w := httptest.NewRecorder()

	middlewareChain := EnsureAutheticatedMiddleware(sharedAPIKeyMiddleware("valid-key")(next))
	middlewareChain.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestSharedAPIKeyMiddleware_EmptyAPIKey(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()

	middlewareChain := EnsureAutheticatedMiddleware(sharedAPIKeyMiddleware("valid-key")(next))
	middlewareChain.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestSharedAPIKeyMiddleware_InvalidBearerFormat(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "valid-key")
	w := httptest.NewRecorder()

	middlewareChain := EnsureAutheticatedMiddleware(sharedAPIKeyMiddleware("valid-key")(next))
	middlewareChain.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestNewMiddleware_DefaultBackend(t *testing.T) {
	cfg := &config.Config{
		AuthBackend:  "",
		APIKeySecret: "valid-key",
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := NewUserIDMiddleware(cfg)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer valid-key")
	w := httptest.NewRecorder()

	middleware(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
