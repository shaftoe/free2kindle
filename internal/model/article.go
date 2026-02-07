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
	Author             string     `json:"author" dynamodbav:"author"`
	Content            string     `json:"content" dynamodbav:"content"`
	HTML               string     `json:"html" dynamodbav:"html"`
	Excerpt            string     `json:"excerpt" dynamodbav:"excerpt"`
	URL                string     `json:"url" dynamodbav:"url"`
	ImageURL           string     `json:"imageUrl" dynamodbav:"imageUrl"`
	PublishedAt        *time.Time `json:"publishedAt,omitempty" dynamodbav:"publishedAt"`
	ExtractedAt        time.Time  `json:"extractedAt" dynamodbav:"extractedAt"`
	WordCount          int        `json:"wordCount" dynamodbav:"wordCount"`
	ReadingTimeMinutes int        `json:"readingTimeMinutes" dynamodbav:"readingTimeMinutes"`
	SourceDomain       string     `json:"sourceDomain" dynamodbav:"sourceDomain"`
	SiteName           string     `json:"siteName" dynamodbav:"siteName"`
	ContentType        string     `json:"contentType" dynamodbav:"contentType"`
	Language           string     `json:"language" dynamodbav:"language"`

	DeliveryStatus       Status     `json:"deliveryStatus" dynamodbav:"deliveryStatus"`
	DeliveryAttemptCount int        `json:"deliveryAttemptCount" dynamodbav:"deliveryAttemptCount"`
	LastDeliveryAttempt  *time.Time `json:"lastDeliveryAttempt,omitempty" dynamodbav:"lastDeliveryAttempt"`
	DeliveryError        string     `json:"deliveryError,omitempty" dynamodbav:"deliveryError"`

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}
