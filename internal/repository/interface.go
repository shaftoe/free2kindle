package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shaftoe/savetoink/internal/model"
)

// Repository defines the interface for article persistence.
type Repository interface {
	Store(ctx context.Context, article *model.Article) error
	GetByAccountAndID(ctx context.Context, account, id string) (*model.Article, error)
	GetMetadataByAccount(
		ctx context.Context,
		account string,
		page, pageSize int,
	) ([]*model.Article, map[string]types.AttributeValue, int, error)
	DeleteByAccountAndID(ctx context.Context, account, id string) error
	DeleteByAccount(ctx context.Context, account string) (int, error)
}
