package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email/mailjet"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/service"
)

func handleCreateArticle(ctx context.Context, req events.LambdaFunctionURLRequest) (*events.APIGatewayProxyResponse, error) {
	apiKey := req.Headers[strings.ToLower(apiKeyHeader)] // NOTICE: for some reason it gets lowered by the Lambda environment

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

	slog.SetDefault(slog.Default().With("url", articleReq.URL))

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
		return nil, err
	}

	body, _ := json.Marshal(ArticleResponse{
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
