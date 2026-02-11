// Package server provides HTTP handlers and middleware for the savetoink application.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shaftoe/savetoink/internal/auth"
	"github.com/shaftoe/savetoink/internal/config"
	"github.com/shaftoe/savetoink/internal/model"
	"github.com/shaftoe/savetoink/internal/service"
)

func newHandlers(
	cfg *config.Config,
	svc service.Interface,
) *handlers {
	return &handlers{
		cfg:     cfg,
		service: svc,
	}
}

func (h *handlers) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}

// handleCreateArticle handles the creation of a new article.
// It delegates all business logic to the service layer.
func (h *handlers) handleCreateArticle(w http.ResponseWriter, r *http.Request) {
	var req articleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{Error: "failed to decode request body: " + err.Error()})
		return
	}

	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{Error: "missing URL in request body"})
		return
	}

	result, err := h.service.CreateArticle(r.Context(), req.URL, auth.GetAccountID(r.Context()))
	if err != nil {
		addLogAttr(r.Context(), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	addLogAttr(r.Context(), slog.String("article_id", result.Article.ID))
	addLogAttr(r.Context(), slog.String("article_url", result.Article.URL))
	addLogAttr(r.Context(), slog.String("message", result.Message))

	if result.Article.DeliveryStatus != "" {
		addLogAttr(r.Context(), slog.String("delivery_status", string(result.Article.DeliveryStatus)))
		addLogAttr(r.Context(), slog.String("email_provider", string(h.cfg.EmailProvider)))
	}

	if result.EmailResp != nil && result.EmailResp.EmailUUID != "" {
		addLogAttr(r.Context(), slog.String("email_uuid", result.EmailResp.EmailUUID))
	}

	if dbErr := h.service.GetDBError(); dbErr != nil {
		addLogAttr(r.Context(), slog.String("db_error", dbErr.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(articleResponse{
		ID:             result.Article.ID,
		Title:          result.Article.Title,
		URL:            result.Article.URL,
		Message:        result.Message,
		DeliveryStatus: string(result.Article.DeliveryStatus),
	})
}
