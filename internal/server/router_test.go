package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shaftoe/savetoink/internal/auth"
	"github.com/shaftoe/savetoink/internal/config"
	"github.com/shaftoe/savetoink/internal/model"
	"github.com/shaftoe/savetoink/internal/service"
)

const (
	testAPIKey           = "test-api-key"
	testArticleTitle     = "Test Article"
	testArticleURL       = "https://example.com/article"
	testOrigin           = "https://example.com"
	statusOk             = "ok"
	statusCreatedMessage = "article sent to Kindle successfully"
	statusEmailDisabled  = "article processed successfully (email sending disabled)"
	methodsAllowed       = "POST, GET, OPTIONS"
)

func createTestRouterWithHandler(h *handlers, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(auth.NewUserIDMiddleware(cfg))
	r.Use(requestIDMiddleware)
	r.Use(corsMiddleware)
	r.Use(jsonContentTypeMiddleware)
	r.Use(loggingMiddleware)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/articles", func(r chi.Router) {
			r.Use(auth.EnsureAutheticatedMiddleware)
			r.Post("/", h.handleCreateArticle)
		})
	})
	return r
}

func createTestHandlerWithMock(
	cfg *config.Config,
) *handlers {
	svc := newMockService(func(_ context.Context, _ string) (*service.ProcessResult, error) {
		return service.NewProcessResult(
			&model.Article{Title: testArticleTitle},
			[]byte("epub data"),
			testArticleURL,
		), nil
	})
	return newHandlers(cfg, svc, nil)
}

func TestNewRouter_RouteRegistration(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: "test-key",
		SendEnabled:  false,
	}
	r := NewRouter(cfg)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		body           []byte
	}{
		{
			name:           "health endpoint - GET",
			method:         "GET",
			path:           "/v1/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "health endpoint - OPTIONS",
			method:         "OPTIONS",
			path:           "/v1/health",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "articles endpoint - POST without auth",
			method:         "POST",
			path:           "/v1/articles",
			expectedStatus: http.StatusUnauthorized,
			body:           []byte(`{"url":"https://example.com"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(tt.body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestNewRouter_404Handler(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: "test-key",
	}
	r := NewRouter(cfg)

	req := httptest.NewRequest("GET", "/v1/unknown", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "not_found" {
		t.Errorf("expected message 'not_found', got '%s'", resp.Error)
	}
}

func TestNewRouter_405Handler(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: "test-key",
	}
	r := NewRouter(cfg)

	req := httptest.NewRequest("PUT", "/v1/health", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "method_not_allowed" {
		t.Errorf("expected message 'method_not_allowed', got '%s'", resp.Error)
	}
}

func TestNewRouter_MiddlewareChain(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: "test-key",
		SendEnabled:  false,
	}

	h := createTestHandlerWithMock(cfg)

	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Use(requestIDMiddleware)
	r.Get("/test", h.handleHealth)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("X-Request-ID") == "" {
		t.Errorf("expected X-Request-ID header to be set")
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected Access-Control-Allow-Origin header to be set")
	}
}

func TestHealthCheckFlow(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: "test-key",
	}
	r := NewRouter(cfg)

	req := httptest.NewRequest("GET", "/v1/health", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp healthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != statusOk {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
}

func TestArticleCreationFlow_Authenticated(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: testAPIKey,
		SendEnabled:  true,
	}

	h := createTestHandlerWithMock(cfg)
	r := createTestRouterWithHandler(h, cfg)

	body := articleRequest{URL: testArticleURL}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testAPIKey)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp articleResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Title != testArticleTitle {
		t.Errorf("expected title '%s', got '%s'", testArticleTitle, resp.Title)
	}
	if resp.Message != statusCreatedMessage {
		t.Errorf("expected message '%s', got '%s'", statusCreatedMessage, resp.Message)
	}
}

func TestArticleCreationFlow_Unauthenticated(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: testAPIKey,
		SendEnabled:  false,
	}

	h := createTestHandlerWithMock(cfg)
	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/articles", func(r chi.Router) {
			r.Use(auth.EnsureAutheticatedMiddleware)
			r.Use(auth.NewUserIDMiddleware(cfg))
			r.Post("/", h.handleCreateArticle)
		})
	})

	body := articleRequest{URL: testArticleURL}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestArticleCreationFlow_EmailDisabled(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: testAPIKey,
		SendEnabled:  false,
	}

	h := createTestHandlerWithMock(cfg)
	r := createTestRouterWithHandler(h, cfg)

	body := articleRequest{URL: testArticleURL}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testAPIKey)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp articleResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != statusEmailDisabled {
		t.Errorf("expected message '%s', got '%s'", statusEmailDisabled, resp.Message)
	}
}

func TestCORS_PreflightFlow(t *testing.T) {
	cfg := &config.Config{
		APIKeySecret: "test-key",
	}
	r := NewRouter(cfg)

	req := httptest.NewRequest("OPTIONS", "/v1/health", http.NoBody)
	req.Header.Set("origin", testOrigin)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != testOrigin {
		t.Errorf("expected Access-Control-Allow-Origin '%s', got '%s'",
			testOrigin, w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Methods") != methodsAllowed {
		t.Errorf("expected Access-Control-Allow-Methods '%s', got '%s'",
			methodsAllowed, w.Header().Get("Access-Control-Allow-Methods"))
	}
}
