package server

import (
	"context"
	"log/slog"
)

type articleRequest struct {
	URL string `json:"url"`
}

type articleResponse struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type handlerDeps struct {
	kindleEmail      string
	senderEmail      string
	mailjetAPIKey    string
	mailjetAPISecret string
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
