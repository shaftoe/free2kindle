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
	Title              string     `json:"title" dynamodbav:"title"`
	Author             string     `json:"author,omitempty" dynamodbav:"author,omitempty"`
	Content            string     `json:"content" dynamodbav:"content"`
	Excerpt            string     `json:"excerpt,omitempty" dynamodbav:"excerpt,omitempty"`
	URL                string     `json:"url" dynamodbav:"url"`
	ImageURL           string     `json:"imageUrl,omitempty" dynamodbav:"imageUrl,omitempty"`
	PublishedAt        *time.Time `json:"publishedAt,omitempty" dynamodbav:"publishedAt,omitempty"`
	WordCount          int        `json:"wordCount" dynamodbav:"wordCount"`
	ReadingTimeMinutes int        `json:"readingTimeMinutes" dynamodbav:"readingTimeMinutes"`
	SourceDomain       string     `json:"sourceDomain" dynamodbav:"sourceDomain"`
	SiteName           string     `json:"siteName" dynamodbav:"siteName"`
	ContentType        string     `json:"contentType,omitempty" dynamodbav:"contentType,omitempty"`
	Language           string     `json:"language,omitempty" dynamodbav:"language,omitempty"`

	LastDeliveryAttemptAt *time.Time `json:"lastDeliveryAttemptAt,omitempty" dynamodbav:"lastDeliveryAttemptAt,omitempty"`
	DeliveryStatus        Status     `json:"deliveryStatus,omitempty" dynamodbav:"deliveryStatus,omitempty"`
	DeliveryError         *string    `json:"deliveryError,omitempty" dynamodbav:"deliveryError,omitempty"`

	CreatedAt time.Time  `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" dynamodbav:"updatedAt,omitempty"`
}
