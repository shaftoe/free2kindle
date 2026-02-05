// Package main implements the AWS Lambda handler for the free2kindle service.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/config"
)

const (
	apiKeyHeader = "X-API-Key"
	version      = "0.0.0-devel"
)

var (
	err error
	cfg *config.Config
)

// setupLogging initializes the logging system for wide events logging
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

func handleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (*events.APIGatewayProxyResponse, error) {
	var resp *events.APIGatewayProxyResponse

	setupLogging(ctx, &req)

	if req.RequestContext.HTTP.Method == http.MethodOptions {
		slog.Debug("options request")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusNoContent,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, " + apiKeyHeader,
			},
		}, nil
	}

	switch req.RequestContext.HTTP.Method + req.RequestContext.HTTP.Path {
	case http.MethodGet + "/api/v1/health":
		slog.Debug("health check")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "ok",
		}, nil

	case http.MethodPost + "/api/v1/articles":
		resp, err = handleCreateArticle(ctx, req)

		if resp != nil {
			slog.Info("article sent")
			return resp, nil
		}
	}

	if err != nil {
		return respondError(http.StatusInternalServerError, "internal_server_error", err.Error())
	}

	slog.Debug("not found")

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       "not found",
	}, nil
}

func respondError(status int, errorType string, message string) (*events.APIGatewayProxyResponse, error) {
	slog.With("status", status).With("error_type", errorType).Info("request failed", "error", message)

	body, _ := json.Marshal(ErrorResponse{
		Error:   errorType,
		Message: message,
	})
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
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

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	cfg, err = config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	lambda.Start(handleRequest)
}
