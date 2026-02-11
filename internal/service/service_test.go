package service

import (
	"context"
	"testing"
	"time"

	"github.com/shaftoe/savetoink/internal/constant"
	"github.com/shaftoe/savetoink/internal/model"
)

type MockRepository struct {
	articles []*model.Article
}

func (m *MockRepository) Store(_ context.Context, article *model.Article) error {
	m.articles = append(m.articles, article)
	return nil
}

func (m *MockRepository) GetByAccountAndID(_ context.Context, account, id string) (*model.Article, error) {
	for _, article := range m.articles {
		if article.Account == account && article.ID == id {
			return article, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) GetByAccount(_ context.Context, account string) ([]*model.Article, error) {
	var result []*model.Article
	for _, article := range m.articles {
		if article.Account == account {
			result = append(result, article)
		}
	}
	return result, nil
}

func (m *MockRepository) DeleteByAccountAndID(_ context.Context, account, id string) error {
	for i, article := range m.articles {
		if article.Account == account && article.ID == id {
			m.articles = append(m.articles[:i], m.articles[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockRepository) DeleteByAccount(_ context.Context, account string) (int, error) {
	initialLen := len(m.articles)
	var filtered []*model.Article
	for _, article := range m.articles {
		if article.Account != account {
			filtered = append(filtered, article)
		}
	}
	m.articles = filtered
	return initialLen - len(m.articles), nil
}

func TestGetArticles(t *testing.T) {
	articles := []*model.Article{
		{Account: "user1", ID: "1", Title: "Article 1", URL: "https://example.com/1", CreatedAt: time.Now()},
		{Account: "user1", ID: "2", Title: "Article 2", URL: "https://example.com/2", CreatedAt: time.Now()},
		{Account: "user1", ID: "3", Title: "Article 3", URL: "https://example.com/3", CreatedAt: time.Now()},
		{Account: "user1", ID: "4", Title: "Article 4", URL: "https://example.com/4", CreatedAt: time.Now()},
		{Account: "user1", ID: "5", Title: "Article 5", URL: "https://example.com/5", CreatedAt: time.Now()},
	}

	tests := []struct {
		name            string
		accountID       string
		page            int
		pageSize        int
		expectedCount   int
		expectedPage    int
		expectedTotal   int
		expectedHasMore bool
	}{
		{
			name:            "first page with page size 2",
			accountID:       "user1",
			page:            1,
			pageSize:        2,
			expectedCount:   2,
			expectedPage:    1,
			expectedTotal:   5,
			expectedHasMore: true,
		},
		{
			name:            "second page with page size 2",
			accountID:       "user1",
			page:            2,
			pageSize:        2,
			expectedCount:   2,
			expectedPage:    2,
			expectedTotal:   5,
			expectedHasMore: true,
		},
		{
			name:            "last page with page size 2",
			accountID:       "user1",
			page:            3,
			pageSize:        2,
			expectedCount:   1,
			expectedPage:    3,
			expectedTotal:   5,
			expectedHasMore: false,
		},
		{
			name:            "page beyond last returns empty",
			accountID:       "user1",
			page:            10,
			pageSize:        2,
			expectedCount:   0,
			expectedPage:    10,
			expectedTotal:   5,
			expectedHasMore: false,
		},
		{
			name:            "get all articles in one page",
			accountID:       "user1",
			page:            1,
			pageSize:        100,
			expectedCount:   5,
			expectedPage:    1,
			expectedTotal:   5,
			expectedHasMore: false,
		},
		{
			name:            "no articles for account",
			accountID:       "user2",
			page:            1,
			pageSize:        10,
			expectedCount:   0,
			expectedPage:    1,
			expectedTotal:   0,
			expectedHasMore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{articles: articles}
			svc := &Service{repo: mockRepo}

			result, err := svc.GetArticles(context.Background(), tt.accountID, tt.page, tt.pageSize)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result.Articles) != tt.expectedCount {
				t.Errorf("expected %d articles, got %d", tt.expectedCount, len(result.Articles))
			}

			if result.Page != tt.expectedPage {
				t.Errorf("expected page %d, got %d", tt.expectedPage, result.Page)
			}

			if result.PageSize != tt.pageSize {
				t.Errorf("expected page_size %d, got %d", tt.pageSize, result.PageSize)
			}

			if result.Total != tt.expectedTotal {
				t.Errorf("expected total %d, got %d", tt.expectedTotal, result.Total)
			}

			if result.HasMore != tt.expectedHasMore {
				t.Errorf("expected has_more %v, got %v", tt.expectedHasMore, result.HasMore)
			}
		})
	}
}

func TestGetArticlesWithNilRepo(t *testing.T) {
	svc := &Service{repo: nil}

	result, err := svc.GetArticles(context.Background(), "user1", 1, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Articles == nil {
		t.Error("expected articles to be initialized, got nil")
	}

	if len(result.Articles) != 0 {
		t.Errorf("expected 0 articles, got %d", len(result.Articles))
	}

	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
}

func TestGetArticlesWithDeliveryStatus(t *testing.T) {
	articles := []*model.Article{
		{
			Account:        "user1",
			ID:             "1",
			Title:          "Article 1",
			URL:            "https://example.com/1",
			CreatedAt:      time.Now(),
			DeliveryStatus: constant.StatusDelivered,
		},
		{
			Account:        "user1",
			ID:             "2",
			Title:          "Article 2",
			URL:            "https://example.com/2",
			CreatedAt:      time.Now(),
			DeliveryStatus: constant.StatusFailed,
			Error:          "email failed",
		},
	}

	mockRepo := &MockRepository{articles: articles}
	svc := &Service{repo: mockRepo}

	result, err := svc.GetArticles(context.Background(), "user1", 1, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Articles) != 2 {
		t.Errorf("expected 2 articles, got %d", len(result.Articles))
	}

	if result.Articles[0].DeliveryStatus != constant.StatusDelivered {
		t.Errorf("expected delivery status %v, got %v", constant.StatusDelivered, result.Articles[0].DeliveryStatus)
	}

	if result.Articles[1].DeliveryStatus != constant.StatusFailed {
		t.Errorf("expected delivery status %v, got %v", constant.StatusFailed, result.Articles[1].DeliveryStatus)
	}

	if result.Articles[1].Error != "email failed" {
		t.Errorf("expected error 'email failed', got '%s'", result.Articles[1].Error)
	}
}
