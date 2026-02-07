package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/service"
)

func TestHandleHealth(t *testing.T) {
	h := newHandlers(&handlerDeps{})

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	w := httptest.NewRecorder()

	h.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp healthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
}

func TestHandleCreateArticleSuccessWithEmail(t *testing.T) {
	cfg := &config.Config{
		DestEmail:        "test@example.com",
		SenderEmail:      "sender@example.com",
		MailjetAPIKey:    "test-key",
		MailjetAPISecret: "test-secret",
		SendEnabled:      true,
	}
	mockService := func(
		_ context.Context, _ *service.Deps, _ *config.Config,
		_ *service.Options, _ string,
	) (*service.Result, error) {
		return &service.Result{
			Article: &content.Article{Title: "Test Article"},
		}, nil
	}
	h := newHandlers(&handlerDeps{
		cfg:        cfg,
		serviceRun: mockService,
	})

	body := articleRequest{URL: "https://example.com/article"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

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
	if resp.URL != "https://example.com/article" {
		t.Errorf("expected URL 'https://example.com/article', got '%s'", resp.URL)
	}
	if resp.Message != "article sent to Kindle successfully" {
		t.Errorf("expected message 'article sent to Kindle successfully', got '%s'", resp.Message)
	}
}

func TestHandleCreateArticleSuccessWithoutEmail(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	mockService := func(
		_ context.Context, _ *service.Deps, _ *config.Config,
		_ *service.Options, _ string,
	) (*service.Result, error) {
		return &service.Result{
			Article: &content.Article{Title: "Test Article"},
		}, nil
	}
	h := newHandlers(&handlerDeps{
		cfg:        cfg,
		serviceRun: mockService,
	})

	body := articleRequest{URL: "https://example.com/article"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp articleResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Title != "Test Article" {
		t.Errorf("expected title 'Test Article', got '%s'", resp.Title)
	}
	if resp.Message != "article processed successfully (email sending disabled)" {
		t.Errorf("expected message 'article processed successfully (email sending disabled)', got '%s'", resp.Message)
	}
}

func TestHandleCreateArticleInvalidJSON(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	h := newHandlers(&handlerDeps{cfg: cfg})

	req := httptest.NewRequest("POST", "/api/v1/articles", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Invalid request body" {
		t.Errorf("expected message 'Invalid request body', got '%s'", resp.Message)
	}
}

func TestHandleCreateArticleMissingURL(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	h := newHandlers(&handlerDeps{cfg: cfg})

	body := articleRequest{URL: ""}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "URL is required" {
		t.Errorf("expected message 'URL is required', got '%s'", resp.Message)
	}
}

func TestHandleCreateArticleServiceError(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	mockService := func(
		_ context.Context, _ *service.Deps, _ *config.Config,
		_ *service.Options, _ string,
	) (*service.Result, error) {
		return nil, &serviceError{msg: "extraction failed"}
	}
	h := newHandlers(&handlerDeps{
		cfg:        cfg,
		serviceRun: mockService,
	})

	body := articleRequest{URL: "https://example.com/article"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Failed to process article: extraction failed" {
		t.Errorf("expected message 'Failed to process article: extraction failed', got '%s'", resp.Message)
	}
}

func TestHandleCreateArticleNilArticle(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	mockService := func(
		_ context.Context, _ *service.Deps, _ *config.Config,
		_ *service.Options, _ string,
	) (*service.Result, error) {
		return &service.Result{
			Article: nil,
		}, nil
	}
	h := newHandlers(&handlerDeps{
		cfg:        cfg,
		serviceRun: mockService,
	})

	body := articleRequest{URL: "https://example.com/article"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Failed to process article: article is nil" {
		t.Errorf("expected message 'Failed to process article: article is nil', got '%s'", resp.Message)
	}
}

type serviceError struct {
	msg string
}

func (e *serviceError) Error() string {
	return e.msg
}
