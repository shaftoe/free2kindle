package types

import "time"

const (
	StatusPending    = "pending"
	StatusDelivering = "delivering"
	StatusDelivered  = "delivered"
	StatusFailed     = "failed"
)

type Article struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	Author             string    `json:"author"`
	Content            string    `json:"content"`
	HTML               string    `json:"html"`
	Excerpt            string    `json:"excerpt"`
	URL                string    `json:"url"`
	ImageURL           string    `json:"imageUrl"`
	PublishedAt        time.Time `json:"publishedAt,omitempty"`
	ExtractedAt        time.Time `json:"extractedAt"`
	WordCount          int       `json:"wordCount"`
	ReadingTimeMinutes int       `json:"readingTimeMinutes"`
	SourceDomain       string    `json:"sourceDomain"`
	SiteName           string    `json:"siteName"`
	ContentType        string    `json:"contentType"`
	Language           string    `json:"language"`

	DeliveryStatus       string     `json:"deliveryStatus"`
	DeliveryAttemptCount int        `json:"deliveryAttemptCount"`
	LastDeliveryAttempt  *time.Time `json:"lastDeliveryAttempt,omitempty"`
	DeliveryError        string     `json:"deliveryError,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
