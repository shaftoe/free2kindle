package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/email"
	"github.com/shaftoe/free2kindle/internal/model"
	"github.com/shaftoe/free2kindle/internal/repository"
	"github.com/shaftoe/free2kindle/internal/service"
)

type MockService struct {
	processFunc func(context.Context, string) (*service.ProcessResult, error)
	sendFunc    func(context.Context, *service.ProcessResult, string) (*email.SendEmailResponse, error)
	writeFunc   func(*service.ProcessResult, string) error
}

func newMockService(processFunc func(context.Context, string) (*service.ProcessResult, error)) *MockService {
	return &MockService{
		processFunc: processFunc,
	}
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

func TestHandleHealth(t *testing.T) {
	h := newHandlers(nil, nil, nil)

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
	svc := newMockService(func(_ context.Context, _ string) (*service.ProcessResult, error) {
		return service.NewProcessResult(
			&model.Article{Title: testArticleTitle},
			[]byte("epub data"),
			"https://example.com/article",
		), nil
	})
	h := newHandlers(cfg, svc, nil)

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
	svc := newMockService(func(_ context.Context, _ string) (*service.ProcessResult, error) {
		return service.NewProcessResult(
			&model.Article{Title: testArticleTitle},
			[]byte("epub data"),
			"https://example.com/article",
		), nil
	})
	h := newHandlers(cfg, svc, nil)

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
	if resp.Message != "article processed successfully (email sending disabled)" {
		t.Errorf("expected message 'article processed successfully (email sending disabled)', got '%s'", resp.Message)
	}
}

func TestHandleCreateArticleInvalidJSON(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	h := newHandlers(cfg, nil, nil)

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

	expectedMsg := "failed to decode request body: invalid character 'i' looking for beginning of value"
	if resp.Message != expectedMsg {
		t.Errorf("expected message '%s', got '%s'", expectedMsg, resp.Message)
	}
}

func TestHandleCreateArticleMissingURL(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	h := newHandlers(cfg, nil, nil)

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

	if resp.Message != "missing URL in request body" {
		t.Errorf("expected message 'missing URL in request body', got '%s'", resp.Message)
	}
}

func TestHandleCreateArticleServiceError(t *testing.T) {
	cfg := &config.Config{
		SendEnabled: false,
	}
	svc := newMockService(func(_ context.Context, _ string) (*service.ProcessResult, error) {
		return nil, &serviceError{msg: "extraction failed"}
	})
	h := newHandlers(cfg, svc, nil)

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
	svc := newMockService(func(_ context.Context, _ string) (*service.ProcessResult, error) {
		return service.NewProcessResult(nil, []byte("epub data"), "https://example.com/article"), nil
	})
	h := newHandlers(cfg, svc, nil)

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

func TestGetArticleIDandCleanURL(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid URL",
			body:    `{"url": "https://example.com/article"}`,
			wantErr: false,
		},
		{
			name:        "invalid JSON",
			body:        `invalid json`,
			wantErr:     true,
			errContains: "failed to decode request body",
		},
		{
			name:        "missing URL",
			body:        `{"url": ""}`,
			wantErr:     true,
			errContains: "missing URL in request body",
		},
		{
			name:        "no URL field",
			body:        `{}`,
			wantErr:     true,
			errContains: "missing URL in request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/articles", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			id, url, err := getArticleIDandCleanURL(req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if id == nil {
				t.Error("expected id to be non-nil")
			}
			if url == nil {
				t.Error("expected url to be non-nil")
			}
		})
	}
}

func TestProcessDBArticleUpdates(t *testing.T) {
	tests := []struct {
		name       string
		repository repository.Repository
	}{
		{
			name:       "with repository",
			repository: &mockRepository{},
		},
		{
			name:       "without repository",
			repository: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandlers(nil, nil, tt.repository)
			ctx := context.Background()

			eg, articlesChan := h.processDBArticleUpdates(ctx)

			article := &model.Article{
				ID:  "test-id",
				URL: "https://example.com",
			}

			go func() {
				articlesChan <- article
				close(articlesChan)
			}()

			if err := eg.Wait(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestEnrichArticle(t *testing.T) {
	tests := []struct {
		name          string
		sendEnabled   bool
		senderEmail   string
		destEmail     string
		wantStatus    model.Status
		wantDelivered bool
	}{
		{
			name:          "send enabled",
			sendEnabled:   true,
			senderEmail:   "sender@example.com",
			destEmail:     "dest@example.com",
			wantStatus:    model.StatusDelivered,
			wantDelivered: true,
		},
		{
			name:          "send disabled",
			sendEnabled:   false,
			wantStatus:    "",
			wantDelivered: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				SendEnabled: tt.sendEnabled,
				SenderEmail: tt.senderEmail,
				DestEmail:   tt.destEmail,
			}
			h := newHandlers(cfg, nil, nil)

			id := "test-id"
			article := &model.Article{}
			var emailResp *email.SendEmailResponse
			if tt.wantDelivered {
				emailResp = &email.SendEmailResponse{EmailUUID: "test-uuid"}
			}

			h.enrichArticle(article, &id, emailResp)

			if article.ID != id {
				t.Errorf("expected ID %q, got %q", id, article.ID)
			}

			if article.DeliveryStatus != tt.wantStatus {
				t.Errorf("expected DeliveryStatus %q, got %q", tt.wantStatus, article.DeliveryStatus)
			}

			if tt.wantDelivered {
				assertDeliveredField(t, "DeliveredFrom", tt.senderEmail, article.DeliveredFrom)
				assertDeliveredField(t, "DeliveredTo", tt.destEmail, article.DeliveredTo)
			} else {
				assertNilField(t, "DeliveredFrom", article.DeliveredFrom)
				assertNilField(t, "DeliveredTo", article.DeliveredTo)
			}
		})
	}
}

func TestEnrichLogs(t *testing.T) {
	tests := []struct {
		name           string
		sendEnabled    bool
		emailResp      *email.SendEmailResponse
		wantMessage    string
		deliveryStatus model.Status
	}{
		{
			name:           "send enabled with email response",
			sendEnabled:    true,
			emailResp:      &email.SendEmailResponse{EmailUUID: "test-uuid"},
			wantMessage:    "article sent to Kindle successfully",
			deliveryStatus: model.StatusDelivered,
		},
		{
			name:           "send enabled without email response",
			sendEnabled:    true,
			emailResp:      nil,
			wantMessage:    "article sent to Kindle successfully",
			deliveryStatus: model.StatusDelivered,
		},
		{
			name:           "send disabled",
			sendEnabled:    false,
			emailResp:      nil,
			wantMessage:    "article processed successfully (email sending disabled)",
			deliveryStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				SendEnabled: tt.sendEnabled,
			}
			h := newHandlers(cfg, nil, nil)

			ctx := context.Background()
			article := &model.Article{
				DeliveryStatus: tt.deliveryStatus,
			}

			msg := h.enrichLogs(ctx, article, tt.emailResp)

			if msg == nil {
				t.Fatal("expected msg to be non-nil")
			}
			if *msg != tt.wantMessage {
				t.Errorf("expected message %q, got %q", tt.wantMessage, *msg)
			}
		})
	}
}

type serviceError struct {
	msg string
}

func (e *serviceError) Error() string {
	return e.msg
}

type mockRepository struct{}

func (m *mockRepository) Store(_ context.Context, _ *model.Article) error {
	return nil
}

func (m *mockRepository) GetByAccountAndID(_ context.Context, _, _ string) (*model.Article, error) {
	return nil, nil
}

func (m *mockRepository) GetByAccount(_ context.Context, _ string) ([]*model.Article, error) {
	return nil, nil
}

func (m *mockRepository) DeleteByAccountAndID(_ context.Context, _, _ string) error {
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s != "" && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}

func assertDeliveredField(t *testing.T, name, want string, got *string) {
	t.Helper()
	if got == nil {
		t.Errorf("expected %s to be %q, got nil", name, want)
		return
	}
	if *got != want {
		t.Errorf("expected %s to be %q, got %q", name, want, *got)
	}
}

func assertNilField(t *testing.T, name string, got *string) {
	t.Helper()
	if got != nil {
		t.Errorf("expected %s to be nil, got %q", name, *got)
	}
}
