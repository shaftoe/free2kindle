package server

import (
	"context"
	"log/slog"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/repository"
	"github.com/shaftoe/free2kindle/internal/service"
)

type articleRequest struct {
	URL string `json:"url"`
}

type articleResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	Message        string `json:"message"`
	DeliveryStatus string `json:"delivery_status,omitempty"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type handlers struct {
	cfg        *config.Config
	service    service.Interface
	repository repository.Repository
}

type contextKey string

type logRecord struct {
	*slog.Record
}

const logRecordKey = contextKey("log_record")

func addLogAttr(ctx context.Context, attr slog.Attr) {
	if record, ok := ctx.Value(logRecordKey).(*logRecord); ok {
		record.AddAttrs(attr)
	}
}
