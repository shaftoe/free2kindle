package repository

import (
	"context"

	"github.com/shaftoe/free2kindle/internal/model"
)

// Repository defines the interface for article persistence.
// Implementations can use different backends (DynamoDB, PostgreSQL, etc.).
type Repository interface {
	Store(ctx context.Context, article *model.Article) error
	GetByAccountAndID(ctx context.Context, account, id string) (*model.Article, error)
	GetByAccount(ctx context.Context, account string) ([]*model.Article, error)
	DeleteByAccountAndID(ctx context.Context, account, id string) error
}
