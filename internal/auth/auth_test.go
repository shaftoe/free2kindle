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

	middleware := NewMiddleware(cfg)

	if middleware == nil {
		t.Fatal("expected middleware to not be nil")
	}
}

func TestSharedAPIKeyMiddleware_ValidAPIKey(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-API-Key", "valid-key")
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

	sharedAPIKeyMiddleware("valid-key")(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestSharedAPIKeyMiddleware_WrongAPIKey(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-API-Key", "wrong-key")
	w := httptest.NewRecorder()

	sharedAPIKeyMiddleware("valid-key")(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestSharedAPIKeyMiddleware_EmptyAPIKey(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-API-Key", "")
	w := httptest.NewRecorder()

	sharedAPIKeyMiddleware("valid-key")(next).ServeHTTP(w, req)

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

	middleware := NewMiddleware(cfg)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()

	middleware(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
