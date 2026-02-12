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
	createFunc          func(context.Context, string, string) (*service.CreateArticleResult, error)
	processFunc         func(context.Context, string) (*service.ProcessResult, error)
	sendFunc            func(context.Context, *service.ProcessResult, string) (*email.SendEmailResponse, error)
	writeFunc           func(*service.ProcessResult, string) error
	getArticlesMetadata func(context.Context, string, int, int) (*service.GetArticlesResult, error)
	deleteArticle       func(context.Context, string, string) (*service.DeleteArticleResult, error)
	deleteAllArticles   func(context.Context, string) (*service.DeleteArticleResult, error)
	dbError             error
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

func (m *MockService) GetArticlesMetadata(
	ctx context.Context,
	accountID string,
	page int,
	pageSize int,
) (*service.GetArticlesResult, error) {
	if m.getArticlesMetadata != nil {
		return m.getArticlesMetadata(ctx, accountID, page, pageSize)
	}
	return &service.GetArticlesResult{
		Articles: []*model.Article{},
		Page:     1,
		PageSize: 20,
		Total:    0,
		HasMore:  false,
	}, nil
}

func (m *MockService) DeleteArticle(
	ctx context.Context,
	accountID string,
	articleID string,
) (*service.DeleteArticleResult, error) {
	if m.deleteArticle != nil {
		return m.deleteArticle(ctx, accountID, articleID)
	}
	return &service.DeleteArticleResult{Deleted: 1}, nil
}

func (m *MockService) DeleteAllArticles(
	ctx context.Context,
	accountID string,
) (*service.DeleteArticleResult, error) {
	if m.deleteAllArticles != nil {
		return m.deleteAllArticles(ctx, accountID)
	}
	return &service.DeleteArticleResult{Deleted: 0}, nil
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

func TestHandleGetArticlesSuccess(t *testing.T) {
	cfg := &config.Config{}
	svc := newMockService(nil)
	svc.getArticlesMetadata = func(_ context.Context, _ string, page, pageSize int) (*service.GetArticlesResult, error) {
		articles := []*model.Article{
			{ID: "1", Title: "Article 1", URL: "https://example.com/1"},
			{ID: "2", Title: "Article 2", URL: "https://example.com/2"},
		}
		return &service.GetArticlesResult{
			Articles: articles,
			Page:     page,
			PageSize: pageSize,
			Total:    5,
			HasMore:  true,
		}, nil
	}
	h := newHandlers(cfg, svc)

	req := httptest.NewRequest("GET", "/v1/articles?page=1&page_size=2", http.NoBody)
	w := httptest.NewRecorder()

	h.handleGetArticles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp listArticlesResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Articles) != 2 {
		t.Errorf("expected 2 articles, got %d", len(resp.Articles))
	}
	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PageSize != 2 {
		t.Errorf("expected page_size 2, got %d", resp.PageSize)
	}
	if resp.Total != 5 {
		t.Errorf("expected total 5, got %d", resp.Total)
	}
	if !resp.HasMore {
		t.Errorf("expected has_more true, got false")
	}
}

func TestHandleGetArticlesDefaultParams(t *testing.T) {
	cfg := &config.Config{}
	svc := newMockService(nil)
	svc.getArticlesMetadata = func(_ context.Context, _ string, page, pageSize int) (*service.GetArticlesResult, error) {
		return &service.GetArticlesResult{
			Articles: []*model.Article{},
			Page:     page,
			PageSize: pageSize,
			Total:    0,
			HasMore:  false,
		}, nil
	}
	h := newHandlers(cfg, svc)

	req := httptest.NewRequest("GET", "/v1/articles", http.NoBody)
	w := httptest.NewRecorder()

	h.handleGetArticles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp listArticlesResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected default page 1, got %d", resp.Page)
	}
	if resp.PageSize != 20 {
		t.Errorf("expected default page_size 20, got %d", resp.PageSize)
	}
}

func TestHandleGetArticlesInvalidParams(t *testing.T) {
	cfg := &config.Config{}
	svc := newMockService(nil)
	svc.getArticlesMetadata = func(_ context.Context, _ string, page, pageSize int) (*service.GetArticlesResult, error) {
		return &service.GetArticlesResult{
			Articles: []*model.Article{},
			Page:     page,
			PageSize: pageSize,
			Total:    0,
			HasMore:  false,
		}, nil
	}
	h := newHandlers(cfg, svc)

	testCases := []struct {
		name         string
		url          string
		expectedPage int
		expectedSize int
	}{
		{"invalid page uses default", "/v1/articles?page=abc&page_size=10", 1, 10},
		{"negative page uses default", "/v1/articles?page=-1&page_size=10", 1, 10},
		{"zero page uses default", "/v1/articles?page=0&page_size=10", 1, 10},
		{"invalid size uses default", "/v1/articles?page=1&page_size=abc", 1, 20},
		{"size too small uses default", "/v1/articles?page=1&page_size=0", 1, 20},
		{"size too large uses max", "/v1/articles?page=1&page_size=200", 1, 100},
		{"negative size uses default", "/v1/articles?page=1&page_size=-10", 1, 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, http.NoBody)
			w := httptest.NewRecorder()

			h.handleGetArticles(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			}

			var resp listArticlesResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Page != tc.expectedPage {
				t.Errorf("expected page %d, got %d", tc.expectedPage, resp.Page)
			}
			if resp.PageSize != tc.expectedSize {
				t.Errorf("expected page_size %d, got %d", tc.expectedSize, resp.PageSize)
			}
		})
	}
}

func TestHandleGetArticlesServiceError(t *testing.T) {
	cfg := &config.Config{}
	svc := newMockService(nil)
	svc.getArticlesMetadata = func(_ context.Context, _ string, _ int, _ int) (*service.GetArticlesResult, error) {
		return nil, &serviceError{msg: "database error"}
	}
	h := newHandlers(cfg, svc)

	req := httptest.NewRequest("GET", "/v1/articles", http.NoBody)
	w := httptest.NewRecorder()

	h.handleGetArticles(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "database error" {
		t.Errorf("expected error 'database error', got '%s'", resp.Error)
	}
}

func TestHandleDeleteArticle(t *testing.T) {
	tests := []struct {
		name         string
		deleteResult *service.DeleteArticleResult
		deleteErr    error
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "success",
			deleteResult: &service.DeleteArticleResult{Deleted: 1},
			deleteErr:    nil,
			expectedCode: http.StatusOK,
			expectedErr:  "",
		},
		{
			name:         "not found",
			deleteResult: &service.DeleteArticleResult{Deleted: 0},
			deleteErr:    nil,
			expectedCode: http.StatusOK,
			expectedErr:  "",
		},
		{
			name:         "service error",
			deleteResult: nil,
			deleteErr:    &serviceError{msg: testDatabaseError},
			expectedCode: http.StatusInternalServerError,
			expectedErr:  testDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			svc := newMockService(nil)
			svc.deleteArticle = func(_ context.Context, _, _ string) (*service.DeleteArticleResult, error) {
				return tt.deleteResult, tt.deleteErr
			}
			h := newHandlers(cfg, svc)

			req := httptest.NewRequest("DELETE", "/v1/articles/123", http.NoBody)
			w := httptest.NewRecorder()

			h.handleDeleteArticle(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.deleteErr != nil {
				var resp model.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Error != tt.expectedErr {
					t.Errorf("expected error '%s', got '%s'", tt.expectedErr, resp.Error)
				}
				return
			}

			var resp deleteArticleResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.Deleted != tt.deleteResult.Deleted {
				t.Errorf("expected deleted %d, got %d", tt.deleteResult.Deleted, resp.Deleted)
			}
		})
	}
}

func TestHandleDeleteAllArticles(t *testing.T) {
	tests := []struct {
		name         string
		deleteResult *service.DeleteArticleResult
		deleteErr    error
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "success",
			deleteResult: &service.DeleteArticleResult{Deleted: 5},
			deleteErr:    nil,
			expectedCode: http.StatusOK,
			expectedErr:  "",
		},
		{
			name:         "service error",
			deleteResult: nil,
			deleteErr:    &serviceError{msg: testDatabaseError},
			expectedCode: http.StatusInternalServerError,
			expectedErr:  testDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			svc := newMockService(nil)
			svc.deleteAllArticles = func(_ context.Context, _ string) (*service.DeleteArticleResult, error) {
				return tt.deleteResult, tt.deleteErr
			}
			h := newHandlers(cfg, svc)

			req := httptest.NewRequest("DELETE", "/v1/articles", http.NoBody)
			w := httptest.NewRecorder()

			h.handleDeleteAllArticles(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.deleteErr != nil {
				var resp model.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Error != tt.expectedErr {
					t.Errorf("expected error '%s', got '%s'", tt.expectedErr, resp.Error)
				}
				return
			}

			var resp deleteArticleResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.Deleted != tt.deleteResult.Deleted {
				t.Errorf("expected deleted %d, got %d", tt.deleteResult.Deleted, resp.Deleted)
			}
		})
	}
}
