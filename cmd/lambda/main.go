// Package main implements the AWS Lambda handler for the free2kindle service.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/config"
)

const (
	apiKeyHeader      = "X-API-Key" // #nosec G101 - this is a header name, not a credential
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"

	version = "0.0.0-devel"
)

var (
	err error
	cfg *config.Config
)

func getCORSHeaders(req *events.LambdaFunctionURLRequest) map[string]string {
	origin := req.Headers["origin"]
	if origin == "" {
		origin = "*"
	}
	return map[string]string{
		contentTypeHeader:                  contentTypeJSON,
		"Access-Control-Allow-Origin":      origin,
		"Access-Control-Allow-Headers":     fmt.Sprintf("%s, %s", contentTypeHeader, apiKeyHeader),
		"Access-Control-Allow-Methods":     fmt.Sprintf("%s, %s, %s", http.MethodPost, http.MethodGet, http.MethodOptions),
		"Access-Control-Allow-Credentials": "true",
	}
}

type ArticleRequest struct {
	URL string `json:"url"`
}

type ArticleResponse struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

// setupLogging initializes the logging system for wide events logging.
// Ref: https://loggingsucks.com/
func setupLogging(ctx context.Context, req *events.LambdaFunctionURLRequest) {
	leveler := slog.LevelInfo
	if cfg.Debug {
		leveler = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: leveler,
	})

	deadline, _ := ctx.Deadline()
	lc, _ := lambdacontext.FromContext(ctx)
	logger := slog.New(handler).
		With("deadline", deadline).
		With("log_level", leveler.String()).
		With("request_id", lc.AwsRequestID).
		With("request_client_ip", req.RequestContext.HTTP.SourceIP).
		With("request_method", req.RequestContext.HTTP.Method).
		With("request_path", req.RequestContext.HTTP.Path).
		With("version", version)

	slog.SetDefault(logger)
}

func handleRequest(ctx context.Context, req *events.LambdaFunctionURLRequest) *events.LambdaFunctionURLResponse {
	setupLogging(ctx, req)

	if req.RequestContext.HTTP.Method == http.MethodOptions {
		return response(req, &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusNoContent,
		})
	}

	switch req.RequestContext.HTTP.Method + req.RequestContext.HTTP.Path {
	case http.MethodGet + "/api/v1/health":
		return response(req, &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusOK,
			Body:       `{"status": "ok"}`,
		})

	case http.MethodPost + "/api/v1/articles":
		return response(req, handleCreateArticle(ctx, req))

	default:
		return response(req, &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"status": "not_found"}`,
		})
	}
}

func response(req *events.LambdaFunctionURLRequest, resp *events.LambdaFunctionURLResponse) *events.LambdaFunctionURLResponse {
	logger := slog.With("status", resp.StatusCode)

	if resp.StatusCode >= http.StatusNoContent {
		logger.Warn("request failed")
	} else {
		logger.Info("request succeeded")
	}

	resp.Headers = getCORSHeaders(req)

	return resp
}

func main() {
	cfg, err = config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	lambda.Start(handleRequest)
}
