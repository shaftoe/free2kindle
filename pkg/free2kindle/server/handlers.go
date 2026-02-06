// Package server provides HTTP handlers and middleware for the free2kindle application.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email/mailjet"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/service"
)

func newHandlers(deps *handlerDeps) *handlers {
	return &handlers{deps: deps}
}

func (h *handlers) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}

func (h *handlers) handleCreateArticle(w http.ResponseWriter, r *http.Request) {
	var req articleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Invalid request body"})
		return
	}

	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "URL is required"})
		return
	}

	addLogAttr(r.Context(), slog.String("article_url", req.URL))

	mailjetConfig := &mailjet.Config{
		APIKey:      h.deps.mailjetAPIKey,
		APISecret:   h.deps.mailjetAPISecret,
		SenderEmail: h.deps.senderEmail,
	}

	svcCfg := &service.Config{
		Extractor:    content.NewExtractor(),
		Generator:    epub.NewGenerator(),
		Sender:       mailjet.NewSender(mailjetConfig),
		SendEmail:    true,
		GenerateEPUB: true,
		KindleEmail:  h.deps.kindleEmail,
		SenderEmail:  h.deps.senderEmail,
		Subject:      "",
		OutputPath:   "",
	}

	result, err := service.Run(r.Context(), svcCfg, req.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Failed to process article: " + err.Error()})
		return
	}

	addLogAttr(r.Context(), slog.String("article_title", result.Title))

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(articleResponse{
		Title:   result.Title,
		URL:     req.URL,
		Status:  "completed",
		Message: "article sent to Kindle successfully",
	})
}
