// Package repository provides article persistence implementations.
package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/model"
)

const (
	attributeNameID           = "id"
	attributeNameURL          = "url"
	attributeNameTitle        = "title"
	attributeNameAuthor       = "author"
	attributeNameExcerpt      = "excerpt"
	attributeNameImageURL     = "imageUrl"
	attributeNamePublishedAt  = "publishedAt"
	attributeNameExtractedAt  = "extractedAt"
	attributeNameWordCount    = "wordCount"
	attributeNameReadingTime  = "readingTimeMinutes"
	attributeNameSourceDomain = "sourceDomain"
	attributeNameSiteName     = "siteName"
	attributeNameContentType  = "contentType"
	attributeNameLanguage     = "language"
)

// DynamoDB implements Repository interface using AWS DynamoDB.
type DynamoDB struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDB creates a new DynamoDB repository instance.
func NewDynamoDB(awsConfig *aws.Config, tableName string) *DynamoDB {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	if awsConfig != nil && awsConfig.Region == "" {
		cfg.Region = awsConfig.Region
	}
	return &DynamoDB{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}
}

// Store saves an article to DynamoDB.
func (d *DynamoDB) Store(ctx context.Context, article *model.Article) error {
	now := time.Now()
	if article.UpdatedAt.IsZero() {
		article.UpdatedAt = now
	}

	if article.CreatedAt.IsZero() {
		article.CreatedAt = now
	}

	item, err := attributevalue.MarshalMap(article)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to store article: %w", err)
	}

	return nil
}

// GetByID implements Repository.GetByID.
func (d *DynamoDB) GetByID(ctx context.Context, id string) (*model.Article, error) {
	resp, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			attributeNameID: &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	if resp.Item == nil {
		return nil, ErrNotFound
	}

	var article model.Article
	if unmarshalErr := attributevalue.UnmarshalMap(resp.Item, &article); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal article: %w", unmarshalErr)
	}

	return &article, nil
}

// GetByURL implements Repository.GetByURL.
func (d *DynamoDB) GetByURL(ctx context.Context, url string) (*model.Article, error) {
	id, err := content.ArticleIDFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to generate article ID: %w", err)
	}

	return d.GetByID(ctx, id)
}

// ListRecent implements Repository.ListRecent.
func (d *DynamoDB) ListRecent(ctx context.Context, limit int) ([]*model.Article, error) {
	if limit <= 0 {
		limit = 10
	}

	if limit > math.MaxInt32 {
		limit = math.MaxInt32
	}

	scanInput := &dynamodb.ScanInput{
		TableName:            aws.String(d.tableName),
		Limit:                aws.Int32(int32(limit)), // #nosec G115 -- limit is already checked against MaxInt32
		ScanFilter:           nil,
		ProjectionExpression: aws.String("id, url, title, extractedAt, deliveryStatus"),
	}

	resp, err := d.client.Scan(ctx, scanInput)
	if err != nil {
		return nil, fmt.Errorf("failed to scan articles: %w", err)
	}

	var articles []*model.Article
	for _, item := range resp.Items {
		var article model.Article
		if unmarshalErr := attributevalue.UnmarshalMap(item, &article); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal article: %w", unmarshalErr)
		}
		articles = append(articles, &article)
	}

	return articles, nil
}

// ErrNotFound is returned when an article is not found.
var ErrNotFound = errors.New("article not found")
