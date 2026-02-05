package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email/mailjet"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/service"
)

// handleCreateArticle handles the creation and email delivery of a new article.
func handleCreateArticle(
	ctx context.Context,
	req *events.LambdaFunctionURLRequest,
) *events.LambdaFunctionURLResponse {
	// NOTICE: for some reason header key gets lowered by the Lambda environment
	apiKey, ok := req.Headers[strings.ToLower(apiKeyHeader)]

	if !ok || apiKey == "" {
		return &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"message": "API key required"}`,
		}
	}

	if apiKey != cfg.APIKeySecret {
		return &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"message": "Invalid API key"}`,
		}
	}

	var articleReq ArticleRequest
	if err := json.Unmarshal([]byte(req.Body), &articleReq); err != nil {
		return &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Invalid request body"}`,
		}
	}

	if articleReq.URL == "" {
		return &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "URL is required"}`,
		}
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
		return &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message": "Failed to process article: %v"}`, err),
		}
	}

	body, _ := json.Marshal(ArticleResponse{
		Title:   result.Title,
		URL:     articleReq.URL,
		Status:  "completed",
		Message: "article sent to Kindle successfully",
	})

	return &events.LambdaFunctionURLResponse{
		StatusCode: http.StatusCreated,
		Body:       string(body),
		Headers:    getCORSHeaders(req),
	}
}
