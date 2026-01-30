package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
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

	funcframework.RegisterHTTPFunction("/", handleHTTP)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path

	switch {
	case path == "/api/v1/health" && r.Method == http.MethodGet:
		handleHealth(w, r)
	case path == "/api/v1/articles" && r.Method == http.MethodPost:
		handleCreateArticle(w, r)
	default:
		respondError(w, http.StatusNotFound, "not_found", "Endpoint not found")
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func handleCreateArticle(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized", "API key required")
		return
	}

	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(cfg.APIKeySecret)) != 1 {
		respondError(w, http.StatusUnauthorized, "unauthorized", "Invalid API key")
		return
	}

	var req ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.URL == "" {
		respondError(w, http.StatusBadRequest, "invalid_request", "URL is required")
		return
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

	ctx := r.Context()
	result, err := service.Run(ctx, svcCfg, req.URL)
	if err != nil {
		log.Printf("Failed to process article: %v", err)
		respondError(w, http.StatusInternalServerError, "processing_failed", "Failed to process article")
		return
	}

	respondJSON(w, http.StatusCreated, ArticleResponse{
		ID:      generateID(),
		Title:   result.Title,
		URL:     req.URL,
		Status:  "completed",
		Message: "Article sent to Kindle successfully",
	})
}

func generateID() string {
	return fmt.Sprintf("%x", os.Getpid())
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, errorType string, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}

func main() {
	log.Println("Function started")
}
