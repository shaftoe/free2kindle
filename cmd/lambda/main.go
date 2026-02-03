package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/config"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email/mailjet"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/service"
)

var (
	cfg *config.Config
)

func init() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
}

type APIGatewayRequest struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Path    string            `json:"path"`
	Method  string            `json:"httpMethod"`
}

type APIGatewayResponse struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	path := req.Path
	method := req.HTTPMethod

	if method == "OPTIONS" {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, X-API-Key",
			},
		}, nil
	}

	if path == "/api/v1/health" && method == http.MethodGet {
		return handleHealth()
	}

	if path == "/api/v1/articles" && method == http.MethodPost {
		return handleCreateArticle(ctx, req)
	}

	return respondError(http.StatusNotFound, "not_found", "Endpoint not found")
}

func handleHealth() (*events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(HealthResponse{Status: "ok"})
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func handleCreateArticle(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	apiKey := req.Headers["X-API-Key"]
	if apiKey == "" {
		return respondError(http.StatusUnauthorized, "unauthorized", "API key required")
	}

	if apiKey != cfg.APIKeySecret {
		return respondError(http.StatusUnauthorized, "unauthorized", "Invalid API key")
	}

	var articleReq ArticleRequest
	if err := json.Unmarshal([]byte(req.Body), &articleReq); err != nil {
		return respondError(http.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if articleReq.URL == "" {
		return respondError(http.StatusBadRequest, "invalid_request", "URL is required")
	}

	mailjetConfig := &mailjet.Config{
		APIKey:      cfg.MailjetAPIKey,
		APISecret:   cfg.MailjetAPISecret,
		SenderEmail: cfg.SenderEmail,
	}

	svcCfg := &service.Config{
		Extractor:    content.NewExtractor(),
		Generator:    epub.NewGenerator(),
		Sender:       mailjet.NewSender(mailjetConfig),
		SendEmail:    true,
		GenerateEPUB: true,
		KindleEmail:  cfg.KindleEmail,
		SenderEmail:  cfg.SenderEmail,
		Subject:      "",
		OutputPath:   "",
	}

	result, err := service.Run(ctx, svcCfg, articleReq.URL)
	if err != nil {
		log.Printf("Failed to process article: %v", err)
		return respondError(http.StatusInternalServerError, "processing_failed", "Failed to process article")
	}

	body, _ := json.Marshal(ArticleResponse{
		ID:      generateID(),
		Title:   result.Title,
		URL:     articleReq.URL,
		Status:  "completed",
		Message: "Article sent to Kindle successfully",
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func respondError(status int, errorType string, message string) (*events.APIGatewayProxyResponse, error) {
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

func generateID() string {
	return "1"
}

type ArticleRequest struct {
	URL string `json:"url"`
}

type ArticleResponse struct {
	ID      string `json:"id"`
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
	lambda.Start(handleRequest)
}
