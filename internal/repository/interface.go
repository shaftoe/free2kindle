package repository

import (
	"context"

	"github.com/shaftoe/free2kindle/internal/model"
)

// Repository defines the interface for article persistence.
// Implementations can use different backends (DynamoDB, PostgreSQL, etc.).
type Repository interface {
	Store(ctx context.Context, article *model.Article) error
	GetByID(ctx context.Context, id string) (*model.Article, error)
	GetByURL(ctx context.Context, url string) (*model.Article, error)
	UpdateDeliveryStatus(ctx context.Context, id, status string, attemptCount int, errorMsg string) error
	ListRecent(ctx context.Context, limit int) ([]*model.Article, error)
}
