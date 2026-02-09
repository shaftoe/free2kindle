package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/shaftoe/free2kindle/internal/auth"
)

type responseStatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseStatusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

const (
	requestIDKey = "request_id"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ""

		if lc, ok := lambdacontext.FromContext(r.Context()); ok {
			requestID = lc.AwsRequestID
		}

		if requestID == "" {
			requestID = r.Header.Get("X-Request-ID")
		}

		if requestID == "" {
			requestID = r.Header.Get("x-amzn-request-id")
		}

		if requestID == "" {
			requestID = generateRequestID()
		}

		ctx := context.WithValue(r.Context(), contextKey(requestIDKey), requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID, _ := r.Context().Value(contextKey(requestIDKey)).(string)
		userID := auth.GetAccountID(r.Context())

		level := slog.LevelInfo
		record := slog.NewRecord(time.Now(), level, "request completed", 0)
		record.AddAttrs(
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("user_id", userID),
		)

		ctx := context.WithValue(r.Context(), logRecordKey, &logRecord{Record: &record})

		recorder := &responseStatusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r.WithContext(ctx))

		latency := time.Since(start)
		statusCode := recorder.status

		if statusCode >= http.StatusInternalServerError {
			record.Level = slog.LevelError
		}

		record.AddAttrs(
			slog.Int("status", statusCode),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.String("client_ip", remoteAddr(r)),
			slog.String("user_agent", r.Header.Get("User-Agent")),
		)

		if err := slog.Default().Handler().Handle(ctx, record); err != nil {
			slog.Error("failed to log request", "error", err)
		}
	})
}

func generateRequestID() string {
	return strings.ReplaceAll(time.Now().Format("20060102-150405.000"), ".", "")
}

func remoteAddr(r *http.Request) string {
	if r.RemoteAddr != "" {
		return r.RemoteAddr
	}
	return "-"
}
