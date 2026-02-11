package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shaftoe/savetoink/internal/config"
	"github.com/shaftoe/savetoink/internal/email"
	"github.com/shaftoe/savetoink/internal/model"
	"github.com/shaftoe/savetoink/internal/service"
)

type MockService struct {
	createFunc  func(context.Context, string, string) (*service.CreateArticleResult, error)
	processFunc func(context.Context, string) (*service.ProcessResult, error)
	sendFunc    func(context.Context, *service.ProcessResult, string) (*email.SendEmailResponse, error)
	writeFunc   func(*service.ProcessResult, string) error
	dbError     error
}

func newMockService(
	createFunc func(context.Context, string, string) (*service.CreateArticleResult, error),
) *MockService {
	return &MockService{
		createFunc: createFunc,
	}
}

func (m *MockService) CreateArticle(
	ctx context.Context,
	_ string,
	_ string,
) (*service.CreateArticleResult, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, "", "")
	}
	return nil, nil
}

func (m *MockService) Process(ctx context.Context, url string) (*service.ProcessResult, error) {
	if m.processFunc != nil {
		return m.processFunc(ctx, url)
	}
	return nil, nil
}

func (m *MockService) Send(
	ctx context.Context,
	result *service.ProcessResult,
	subject string,
) (*email.SendEmailResponse, error) {
	if m.sendFunc != nil {
		return m.sendFunc(ctx, result, subject)
	}
	return &email.SendEmailResponse{
		Status:    "success",
		Message:   "sent",
		EmailUUID: "test-uuid",
	}, nil
}

func (m *MockService) WriteToFile(result *service.ProcessResult, outputPath string) error {
	if m.writeFunc != nil {
		return m.writeFunc(result, outputPath)
	}
	return nil
}

func (m *MockService) GetDBError() error {
	return m.dbError
}

func TestHandleHealth(t *testing.T) {
	h := newHandlers(nil, nil)

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
	svc := newMockService(func(_ context.Context, _ string, _ string) (*service.CreateArticleResult, error) {
		return &service.CreateArticleResult{
			Article: &model.Article{
				ID:    "test-id",
				Title: testArticleTitle,
				URL:   "https://example.com/article",
			},
			Message: "article sent to Kindle successfully",
			EmailResp: &email.SendEmailResponse{
				Status:    "success",
				Message:   "sent",
				EmailUUID: "test-uuid",
			},
		}, nil
	})
	h := newHandlers(cfg, svc)

	body := articleRequest{URL: "https://example.com/article"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader(bodyBytes))
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
	svc := newMockService(func(_ context.Context, _ string, _ string) (*service.CreateArticleResult, error) {
		return &service.CreateArticleResult{
			Article: &model.Article{
				ID:    "test-id",
				Title: testArticleTitle,
				URL:   "https://example.com/article",
			},
			Message:   "article processed successfully (email sending disabled)",
			EmailResp: nil,
		}, nil
	})
	h := newHandlers(cfg, svc)

	body := articleRequest{URL: "https://example.com/article"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader(bodyBytes))
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
	if resp.Message != "article processed successfully (email sending disabled)" {
		t.Errorf("expected message 'article processed successfully (email sending disabled)', got '%s'", resp.Message)
	}
}

func TestHandleCreateArticleInvalidJSON(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	h := newHandlers(cfg, nil)

	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedMsg := "failed to decode request body: invalid character 'i' looking for beginning of value"
	if resp.Error != expectedMsg {
		t.Errorf("expected message '%s', got '%s'", expectedMsg, resp.Error)
	}
}

func TestHandleCreateArticleMissingURL(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	h := newHandlers(cfg, nil)

	body := articleRequest{URL: ""}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "missing URL in request body" {
		t.Errorf("expected message 'missing URL in request body', got '%s'", resp.Error)
	}
}

func TestHandleCreateArticleServiceError(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	svc := newMockService(func(_ context.Context, _ string, _ string) (*service.CreateArticleResult, error) {
		return nil, &serviceError{msg: "extraction failed"}
	})
	h := newHandlers(cfg, svc)

	body := articleRequest{URL: "https://example.com/article"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateArticle(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "extraction failed" {
		t.Errorf("expected message 'extraction failed', got '%s'", resp.Error)
	}
}

type serviceError struct {
	msg string
}

func (e *serviceError) Error() string {
	return e.msg
}
