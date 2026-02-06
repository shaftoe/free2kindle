package http

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

const apiKeyHeader = "X-API-Key"
const requestIDKey = "request_id"
const lambdaCtxKey = "lambda_ctx"

type contextKey string

type responseBodyWriter struct {
	http.ResponseWriter
	body   bytes.Buffer
	status int
	wrote  bool
}

func (w *responseBodyWriter) WriteHeader(code int) {
	if !w.wrote {
		w.status = code
		w.wrote = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func (w *responseBodyWriter) Message() string {
	type message struct {
		Message string `json:"message"`
	}

	var msg message
	if err := json.Unmarshal(w.body.Bytes(), &msg); err != nil {
		return ""
	}

	return msg.Message
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(apiKeySecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(apiKeyHeader)
			if apiKey == "" || apiKey != apiKeySecret {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid API key"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
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

		wrapped := &responseBodyWriter{
			ResponseWriter: w,
			body:           bytes.Buffer{},
		}

		requestID, _ := r.Context().Value(contextKey(requestIDKey)).(string)

		var lambdaCtx *lambdacontext.LambdaContext
		if lc, ok := lambdacontext.FromContext(r.Context()); ok {
			lambdaCtx = lc
		}

		next.ServeHTTP(wrapped, r)

		latency := time.Since(start)
		statusCode := wrapped.Status()
		msg := wrapped.Message()

		level := slog.LevelInfo
		if statusCode >= 400 && statusCode < 500 {
			level = slog.LevelWarn
		} else if statusCode >= 500 {
			level = slog.LevelError
		}

		record := slog.NewRecord(time.Now(), level, "request completed", 0)
		record.AddAttrs(
			slog.String("timestamp", time.Now().UTC().Format(time.RFC3339)),
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", statusCode),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.String("client_ip", remoteAddr(r)),
			slog.String("user_agent", r.Header.Get("User-Agent")),
		)

		if msg != "" {
			record.AddAttrs(slog.String("message", msg))
		}

		if lambdaCtx != nil {
			record.AddAttrs(
				slog.String("lambda_request_id", lambdaCtx.AwsRequestID),
				slog.String("lambda_function", lambdaCtx.InvokedFunctionArn),
			)
		}

		if statusCode >= 400 {
			record.AddAttrs(slog.String("response_body", truncateString(wrapped.body.String(), 500)))
		}

		slog.Default().Handler().Handle(r.Context(), record)
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

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
