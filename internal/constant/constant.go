// Package constant provides shared constants used across the free2kindle application.
package constant

import "time"

// RunMode defines the application execution mode.
type RunMode string

const (
	// ModeCLI indicates CLI execution mode.
	ModeCLI RunMode = "cli"
	// ModeServer indicates server execution mode.
	ModeServer RunMode = "server"
)

// AuthBackend defines the authentication backend type.
type AuthBackend string

const (
	// AuthBackendSharedAPIKey indicates shared API key authentication.
	AuthBackendSharedAPIKey AuthBackend = "shared_api_key"
)

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

// HTTP server timeout constants.
const (
	// ReadTimeout is the maximum duration for reading the entire request, including the body.
	ReadTimeout = 5 * time.Second
	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout = 10 * time.Second
	// IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
	IdleTimeout = 15 * time.Second
)

// EmailProvider defines the email provider type.
type EmailProvider string

const (
	// EmailBackendMailjet indicates MailJet email backend.
	EmailBackendMailjet EmailProvider = "mailjet"
)
