package repository

import (
	"context"

	"github.com/shaftoe/free2kindle/internal/types"
)

type Article = types.Article

const (
	StatusPending    = types.StatusPending
	StatusDelivering = types.StatusDelivering
	StatusDelivered  = types.StatusDelivered
	StatusFailed     = types.StatusFailed
)

// Repository defines the interface for article persistence.
// Implementations can use different backends (DynamoDB, PostgreSQL, etc.).
type Repository interface {
	Store(ctx context.Context, article *Article) error
	GetByID(ctx context.Context, id string) (*Article, error)
	GetByURL(ctx context.Context, url string) (*Article, error)
	UpdateDeliveryStatus(ctx context.Context, id, status string, attemptCount int, errorMsg string) error
	ListRecent(ctx context.Context, limit int) ([]*Article, error)
}
