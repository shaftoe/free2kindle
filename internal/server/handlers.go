// Package server provides HTTP handlers and middleware for the free2kindle application.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/email"
	"github.com/shaftoe/free2kindle/internal/email/mailjet"
	"github.com/shaftoe/free2kindle/internal/epub"
	"github.com/shaftoe/free2kindle/internal/model"
	"github.com/shaftoe/free2kindle/internal/service"
	"golang.org/x/sync/errgroup"
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

func getArticleIDandRequest(r *http.Request) (*articleRequest, *string, error) {
	var req *articleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	if req.URL == "" {
		return nil, nil, errors.New("missing URL in request body")
	}

	id, err := content.ArticleIDFromURL(req.URL)
	if err != nil {
		return nil, nil, err
	}

	return req, &id, nil
}

func (h *handlers) processDBArticleUpdates(ctx context.Context) (*errgroup.Group, chan<- *model.Article) {
	eg := &errgroup.Group{}
	articles := make(chan *model.Article)

	eg.Go(func() error {

		if h.deps.repository != nil {
			for article := range articles {
				if storeErr := h.deps.repository.Store(ctx, article); storeErr != nil {
					addLogAttr(ctx, slog.String("db_error", storeErr.Error()))
				}
			}
		}

		return nil
	})

	return eg, articles
}

func (h *handlers) handleCreateArticle(w http.ResponseWriter, r *http.Request) {
	req, id, err := getArticleIDandRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
		return
	}

	addLogAttr(r.Context(), slog.String("article_id", *id))

	eg, artChan := h.processDBArticleUpdates(r.Context())

	artChan <- &model.Article{
		ID:  *id,
		URL: req.URL,
	}

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

	if h.deps.cfg.SendEnabled {
		result.Article.DeliveryStatus = model.StatusDelivered
	}

	result.Article.ID = *id

	addLogAttr(r.Context(), slog.String("article_title", result.Article.Title))
	addLogAttr(r.Context(), slog.String("article_id", result.Article.ID))

	if result.Article.DeliveryStatus != "" {
		addLogAttr(r.Context(), slog.String("delivery_status", string(result.Article.DeliveryStatus)))
	}

	artChan <- result.Article
	close(artChan)

	message := messageSentToKindle
	if !h.deps.cfg.SendEnabled {
		message = messageEmailDisabled
	}

	addLogAttr(r.Context(), slog.String("message", message))

	_ = eg.Wait()

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(articleResponse{
		ID:             result.Article.ID,
		Title:          result.Article.Title,
		URL:            req.URL,
		Message:        message,
		DeliveryStatus: string(result.Article.DeliveryStatus),
	})
}
