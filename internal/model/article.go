// Package model provides data models used throughout the application.
package model

import "time"

// Status represents the delivery status of an article.
type Status string

const (
	// StatusPending indicates that the article is pending delivery.
	StatusPending Status = "pending"
	// StatusDelivered indicates that the article has been successfully delivered.
	StatusDelivered Status = "delivered"
	// StatusFailed indicates that the article delivery has failed.
	StatusFailed Status = "failed"
)

// Article represents all article data including content, metadata, and delivery status.
type Article struct {
	ID                 string     `json:"id" dynamodbav:"id"`
	URL                string     `json:"url" dynamodbav:"url"`
	Title              string     `json:"title" dynamodbav:"title"`
	Content            string     `json:"content" dynamodbav:"content"`
	SiteName           string     `json:"siteName" dynamodbav:"siteName"`
	SourceDomain       string     `json:"sourceDomain" dynamodbav:"sourceDomain"`
	WordCount          int        `json:"wordCount" dynamodbav:"wordCount"`
	ReadingTimeMinutes int        `json:"readingTimeMinutes" dynamodbav:"readingTimeMinutes"`
	CreatedAt          time.Time  `json:"createdAt" dynamodbav:"createdAt"`
	PublishedAt        *time.Time `json:"publishedAt,omitempty" dynamodbav:"publishedAt,omitempty"`

	// optional metadata
	Author      string `json:"author,omitempty" dynamodbav:"author,omitempty"`
	Excerpt     string `json:"excerpt,omitempty" dynamodbav:"excerpt,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty" dynamodbav:"imageUrl,omitempty"`
	ContentType string `json:"contentType,omitempty" dynamodbav:"contentType,omitempty"`
	Language    string `json:"language,omitempty" dynamodbav:"language,omitempty"`

	//// email delivery metadata
	DeliveryStatus     Status  `json:"deliveryStatus,omitempty" dynamodbav:"deliveryStatus,omitempty"`
	DeliveredFrom      *string `json:"deliveredFrom,omitempty" dynamodbav:"deliveredFrom,omitempty"`
	DeliveredTo        *string `json:"deliveredTo,omitempty" dynamodbav:"deliveredTo,omitempty"`
	DeliveredEmailUUID *string `json:"deliveredEmailUUID,omitempty" dynamodbav:"deliveredEmailUUID,omitempty"`
}
