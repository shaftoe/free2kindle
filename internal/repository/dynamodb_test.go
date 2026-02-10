// Package repository provides tests for article persistence implementations.
package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shaftoe/savetoink/internal/constant"
	"github.com/shaftoe/savetoink/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func isResourceNotFound(err error) bool {
	var smithyErr interface {
		ErrorCode() string
	}
	return errors.As(err, &smithyErr) && smithyErr.ErrorCode() == "ResourceNotFoundException"
}

func skipIfTableNotFound(t *testing.T, err error) {
	t.Helper()
	if isResourceNotFound(err) {
		t.Skip("DynamoDB table not found, skipping integration test")
	}
}

func TestDynamoDB_Store(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	article := &model.Article{
		Account:   "test@example.com",
		ID:        "test-id-1",
		URL:       "https://example.com/test",
		Title:     "Test Article",
		Content:   "<p>Test content</p>",
		SiteName:  "Example Site",
		CreatedAt: time.Now(),
	}

	err := repo.Store(ctx, article)
	skipIfTableNotFound(t, err)
	require.NoError(t, err)
}

func TestDynamoDB_Store_RequiresAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	article := &model.Article{
		ID:        "test-id-1",
		URL:       "https://example.com/test",
		Title:     "Test Article",
		Content:   "<p>Test content</p>",
		CreatedAt: time.Now(),
	}

	err := repo.Store(ctx, article)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account field is required")
}

func TestDynamoDB_GetByAccountAndID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	expected := &model.Article{
		Account:   "test@example.com",
		ID:        "test-id-2",
		URL:       "https://example.com/test2",
		Title:     "Test Article 2",
		Content:   "<p>Test content 2</p>",
		SiteName:  "Example Site",
		CreatedAt: time.Now(),
	}

	err := repo.Store(ctx, expected)
	skipIfTableNotFound(t, err)
	require.NoError(t, err)

	actual, err := repo.GetByAccountAndID(ctx, "test@example.com", "test-id-2")
	skipIfTableNotFound(t, err)
	require.NoError(t, err)
	assert.Equal(t, expected.Account, actual.Account)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.URL, actual.URL)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Content, actual.Content)
}

func TestDynamoDB_GetByAccountAndID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	_, err := repo.GetByAccountAndID(ctx, "test@example.com", "non-existent-id")
	skipIfTableNotFound(t, err)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestDynamoDB_GetByAccountAndID_WrongAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	article := &model.Article{
		Account:   "user1@example.com",
		ID:        "test-id-3",
		URL:       "https://example.com/test3",
		Title:     "Test Article 3",
		Content:   "<p>Test content 3</p>",
		CreatedAt: time.Now(),
	}

	err := repo.Store(ctx, article)
	skipIfTableNotFound(t, err)
	require.NoError(t, err)

	_, err = repo.GetByAccountAndID(ctx, "user2@example.com", "test-id-3")
	skipIfTableNotFound(t, err)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestDynamoDB_GetByAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	account := "test@example.com"
	articles := []*model.Article{
		{
			Account:   account,
			ID:        "test-id-4",
			URL:       "https://example.com/test4",
			Title:     "Test Article 4",
			Content:   "<p>Test content 4</p>",
			CreatedAt: time.Now(),
		},
		{
			Account:   account,
			ID:        "test-id-5",
			URL:       "https://example.com/test5",
			Title:     "Test Article 5",
			Content:   "<p>Test content 5</p>",
			CreatedAt: time.Now(),
		},
		{
			Account:   "other@example.com",
			ID:        "test-id-6",
			URL:       "https://example.com/test6",
			Title:     "Test Article 6",
			Content:   "<p>Test content 6</p>",
			CreatedAt: time.Now(),
		},
	}

	for _, article := range articles {
		err := repo.Store(ctx, article)
		skipIfTableNotFound(t, err)
		require.NoError(t, err)
	}

	retrieved, err := repo.GetByAccount(ctx, account)
	skipIfTableNotFound(t, err)
	require.NoError(t, err)
	assert.Len(t, retrieved, 2)

	accountIDs := make(map[string]bool)
	for _, article := range retrieved {
		accountIDs[article.ID] = true
	}
	assert.True(t, accountIDs["test-id-4"])
	assert.True(t, accountIDs["test-id-5"])
	assert.False(t, accountIDs["test-id-6"])
}

func TestDynamoDB_GetByAccount_Empty(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	retrieved, err := repo.GetByAccount(ctx, "non-existent@example.com")
	skipIfTableNotFound(t, err)
	require.NoError(t, err)
	assert.Empty(t, retrieved)
}

func TestDynamoDB_DeleteByAccountAndID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	article := &model.Article{
		Account:   "test@example.com",
		ID:        "test-id-7",
		URL:       "https://example.com/test7",
		Title:     "Test Article 7",
		Content:   "<p>Test content 7</p>",
		CreatedAt: time.Now(),
	}

	err := repo.Store(ctx, article)
	skipIfTableNotFound(t, err)
	require.NoError(t, err)

	err = repo.DeleteByAccountAndID(ctx, "test@example.com", "test-id-7")
	skipIfTableNotFound(t, err)
	require.NoError(t, err)

	_, err = repo.GetByAccountAndID(ctx, "test@example.com", "test-id-7")
	skipIfTableNotFound(t, err)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestDynamoDB_DeleteByAccountAndID_WrongAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	article := &model.Article{
		Account:   "user1@example.com",
		ID:        "test-id-8",
		URL:       "https://example.com/test8",
		Title:     "Test Article 8",
		Content:   "<p>Test content 8</p>",
		CreatedAt: time.Now(),
	}

	err := repo.Store(ctx, article)
	skipIfTableNotFound(t, err)
	require.NoError(t, err)

	err = repo.DeleteByAccountAndID(ctx, "user2@example.com", "test-id-8")
	skipIfTableNotFound(t, err)
	assert.Error(t, err)

	_, err = repo.GetByAccountAndID(ctx, "user1@example.com", "test-id-8")
	skipIfTableNotFound(t, err)
	assert.NoError(t, err)
}

func TestDynamoDB_UpdateArticle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupTestDynamoDB(t)
	ctx := context.Background()

	original := &model.Article{
		Account:   "test@example.com",
		ID:        "test-id-9",
		URL:       "https://example.com/test9",
		Title:     "Test Article 9",
		Content:   "<p>Test content 9</p>",
		CreatedAt: time.Now(),
	}

	err := repo.Store(ctx, original)
	skipIfTableNotFound(t, err)
	require.NoError(t, err)

	updated := &model.Article{
		Account:            "test@example.com",
		ID:                 "test-id-9",
		URL:                "https://example.com/test9",
		Title:              "Updated Article 9",
		Content:            "<p>Updated content 9</p>",
		CreatedAt:          original.CreatedAt,
		DeliveryStatus:     constant.StatusDelivered,
		DeliveredFrom:      stringPtr("sender@example.com"),
		DeliveredTo:        stringPtr("kindle@example.com"),
		DeliveredEmailUUID: stringPtr("email-uuid-123"),
		DeliveredBy:        constant.EmailBackendMailjet,
	}

	err = repo.Store(ctx, updated)
	skipIfTableNotFound(t, err)
	require.NoError(t, err)

	retrieved, err := repo.GetByAccountAndID(ctx, "test@example.com", "test-id-9")
	skipIfTableNotFound(t, err)
	require.NoError(t, err)
	assert.Equal(t, "Updated Article 9", retrieved.Title)
	assert.Equal(t, constant.StatusDelivered, retrieved.DeliveryStatus)
	assert.Equal(t, "sender@example.com", *retrieved.DeliveredFrom)
}

func stringPtr(s string) *string {
	return &s
}

func setupTestDynamoDB(t *testing.T) *DynamoDB {
	t.Helper()

	tableName := "test-savetoink-articles"
	repo := NewDynamoDB(nil, tableName)

	t.Cleanup(func() {
		ctx := context.Background()
		articles, err := repo.GetByAccount(ctx, "test@example.com")
		if err != nil {
			return
		}

		for _, article := range articles {
			_ = repo.DeleteByAccountAndID(ctx, article.Account, article.ID)
		}
	})

	return repo
}
