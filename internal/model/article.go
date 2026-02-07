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
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Author             string     `json:"author"`
	Content            string     `json:"content"`
	HTML               string     `json:"html"`
	Excerpt            string     `json:"excerpt"`
	URL                string     `json:"url"`
	ImageURL           string     `json:"imageUrl"`
	PublishedAt        *time.Time `json:"publishedAt,omitempty"`
	ExtractedAt        time.Time  `json:"extractedAt"`
	WordCount          int        `json:"wordCount"`
	ReadingTimeMinutes int        `json:"readingTimeMinutes"`
	SourceDomain       string     `json:"sourceDomain"`
	SiteName           string     `json:"siteName"`
	ContentType        string     `json:"contentType"`
	Language           string     `json:"language"`

	DeliveryStatus       Status     `json:"deliveryStatus"`
	DeliveryAttemptCount int        `json:"deliveryAttemptCount"`
	LastDeliveryAttempt  *time.Time `json:"lastDeliveryAttempt,omitempty"`
	DeliveryError        string     `json:"deliveryError,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
