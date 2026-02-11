// Package repository provides article persistence implementations.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shaftoe/savetoink/internal/constant"
	"github.com/shaftoe/savetoink/internal/model"
)

const (
	attributeNameAccount = "account"
	attributeNameID      = "id"
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

	if article.Account == "" {
		return errors.New("account field is required")
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

// GetByAccountAndID implements Repository.GetByAccountAndID.
func (d *DynamoDB) GetByAccountAndID(ctx context.Context, account, id string) (*model.Article, error) {
	resp, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			attributeNameAccount: &types.AttributeValueMemberS{Value: account},
			attributeNameID:      &types.AttributeValueMemberS{Value: id},
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

// GetByAccount implements Repository.GetByAccount.
func (d *DynamoDB) GetByAccount(ctx context.Context, account string) ([]*model.Article, error) {
	resp, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.tableName),
		KeyConditionExpression: aws.String("#account = :account"),
		ExpressionAttributeNames: map[string]string{
			"#account": attributeNameAccount,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":account": &types.AttributeValueMemberS{Value: account},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
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

// DeleteByAccountAndID implements Repository.DeleteByAccountAndID.
func (d *DynamoDB) DeleteByAccountAndID(ctx context.Context, account, id string) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			attributeNameAccount: &types.AttributeValueMemberS{Value: account},
			attributeNameID:      &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	return nil
}

// DeleteByAccount implements Repository.DeleteByAccount.
func (d *DynamoDB) DeleteByAccount(ctx context.Context, account string) error {
	articles, err := d.GetByAccount(ctx, account)
	if err != nil {
		return fmt.Errorf("failed to get articles for deletion: %w", err)
	}

	if len(articles) == 0 {
		return nil
	}

	for i := 0; i < len(articles); i += constant.DynamoDBBatchSize {
		end := min(i+constant.DynamoDBBatchSize, len(articles))

		var writeReqs []types.WriteRequest
		for _, article := range articles[i:end] {
			writeReqs = append(writeReqs, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						attributeNameAccount: &types.AttributeValueMemberS{Value: article.Account},
						attributeNameID:      &types.AttributeValueMemberS{Value: article.ID},
					},
				},
			})
		}

		_, err = d.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				d.tableName: writeReqs,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete batch of articles: %w", err)
		}
	}

	return nil
}

// ErrNotFound is returned when an article is not found.
var ErrNotFound = errors.New("article not found")
