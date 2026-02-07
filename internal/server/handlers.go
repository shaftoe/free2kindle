// Package server provides HTTP handlers and middleware for the free2kindle application.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/email"
	"github.com/shaftoe/free2kindle/internal/email/mailjet"
	"github.com/shaftoe/free2kindle/internal/epub"
	"github.com/shaftoe/free2kindle/internal/service"
)

const (
	messageSentToKindle  = "article sent to Kindle successfully"
	messageEmailDisabled = "article processed successfully (email sending disabled)"
)

func newHandlers(deps *handlerDeps) *handlers {
	if deps.serviceRun == nil {
		deps.serviceRun = service.Run
	}
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

	var sender email.Sender
	if h.deps.cfg.SendEnabled {
		sender = mailjet.NewSender(h.deps.cfg.MailjetAPIKey, h.deps.cfg.MailjetAPISecret, h.deps.cfg.SenderEmail)
	}

	d := service.NewDeps(
		content.NewExtractor(),
		epub.NewGenerator(),
		sender,
	)

	opts := service.NewOptions(h.deps.cfg.SendEnabled, true, "", "")

	result, err := h.deps.serviceRun(r.Context(), d, h.deps.cfg, opts, req.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Failed to process article: " + err.Error()})
		return
	}

	if result.Article == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Failed to process article: article is nil"})
		return
	}

	addLogAttr(r.Context(), slog.String("article_title", result.Article.Title))
	addLogAttr(r.Context(), slog.String("article_id", result.Article.ID))

	message := messageSentToKindle
	if !h.deps.cfg.SendEnabled {
		message = messageEmailDisabled
	}

	addLogAttr(r.Context(), slog.String("message", message))

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(articleResponse{
		Title:   result.Article.Title,
		URL:     req.URL,
		Message: message,
	})
}
