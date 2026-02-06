package server

import (
	"context"
	"log/slog"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/config"
)

type articleRequest struct {
	URL string `json:"url"`
}

type articleResponse struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Message string `json:"message"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type handlerDeps struct {
	cfg *config.Config
}

type handlers struct {
	deps *handlerDeps
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
