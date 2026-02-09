// Package server provides HTTP handlers and middleware for the free2kindle application.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/constant"
	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/email"
	"github.com/shaftoe/free2kindle/internal/model"
	"github.com/shaftoe/free2kindle/internal/repository"
	"github.com/shaftoe/free2kindle/internal/service"
	"golang.org/x/sync/errgroup"
)

const (
	messageSentToKindle  = "article sent to Kindle successfully"
	messageEmailDisabled = "article processed successfully (email sending disabled)"
)

func newHandlers(
	cfg *config.Config,
	svc service.Interface,
	repo repository.Repository,
) *handlers {
	return &handlers{
		cfg:        cfg,
		service:    svc,
		repository: repo,
	}
}

func (h *handlers) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}

func getArticleIDandCleanURL(r *http.Request) (id, url *string, err error) {
	var req *articleRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	if req.URL == "" {
		return nil, nil, errors.New("missing URL in request body")
	}

	var cleaned string
	cleaned, err = content.CleanURL(req.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to clean URL: %w", err)
	}
	url = &cleaned

	var articleID string
	articleID, err = content.ArticleIDFromURL(cleaned)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate article ID: %w", err)
	}
	id = &articleID

	return id, url, nil
}

func (h *handlers) processDBArticleUpdates(ctx context.Context) (eg *errgroup.Group, articles chan *model.Article) {
	eg = &errgroup.Group{}
	articles = make(chan *model.Article)

	eg.Go(func() error {
		for article := range articles {
			if h.repository != nil {
				if storeErr := h.repository.Store(ctx, article); storeErr != nil {
					addLogAttr(ctx, slog.String("db_error", storeErr.Error()))
				}
			}
		}

		return nil
	})

	return eg, articles
}

// handleCreateArticle handles the creation of a new article.
// It:
// - downloads/processes the article.
// - updates/stores metadata in the repository.
// - (optionally) sends the article to the Kindle.
func (h *handlers) handleCreateArticle(w http.ResponseWriter, r *http.Request) {
	id, cleanURL, err := getArticleIDandCleanURL(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
		return
	}

	addLogAttr(r.Context(), slog.String("article_id", *id))
	addLogAttr(r.Context(), slog.String("article_url", *cleanURL))

	eg, articlesChan := h.processDBArticleUpdates(r.Context())

	articlesChan <- &model.Article{
		Account:   h.cfg.Account,
		ID:        *id,
		URL:       *cleanURL,
		CreatedAt: time.Now(),
	}

	result, err := h.service.Process(r.Context(), *cleanURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Failed to process article: " + err.Error()})
		return
	}

	if result.Article() == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Failed to process article: article is nil"})
		return
	}

	var emailResp *email.SendEmailResponse
	if h.cfg.SendEnabled {
		emailResp, err = h.service.Send(r.Context(), result, "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "Failed to send email: " + err.Error()})
			return
		}
	}

	h.enrichArticle(result.Article(), id, emailResp)
	articlesChan <- result.Article()
	close(articlesChan)
	msg := h.enrichLogs(r.Context(), result.Article(), emailResp)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(articleResponse{
		ID:             result.Article().ID,
		Title:          result.Article().Title,
		URL:            *cleanURL,
		Message:        *msg,
		DeliveryStatus: string(result.Article().DeliveryStatus),
	})

	_ = eg.Wait()
}

func (h *handlers) enrichArticle(article *model.Article, id *string, emailResp *email.SendEmailResponse) {
	article.Account = h.cfg.Account
	article.ID = *id

	if !h.cfg.SendEnabled {
		return
	}

	if emailResp == nil {
		article.DeliveryStatus = constant.StatusFailed
		return
	}

	article.DeliveryStatus = constant.StatusDelivered
	article.DeliveredFrom = &h.cfg.SenderEmail
	article.DeliveredTo = &h.cfg.DestEmail
	article.DeliveredEmailUUID = &emailResp.EmailUUID
	article.DeliveredBy = h.cfg.EmailProvider
}

func (h *handlers) enrichLogs(
	ctx context.Context,
	article *model.Article,
	emailResp *email.SendEmailResponse,
) (msg *string) {
	message := messageSentToKindle
	if !h.cfg.SendEnabled {
		message = messageEmailDisabled
	}
	addLogAttr(ctx, slog.String("message", message))

	if article.DeliveryStatus != "" {
		addLogAttr(ctx, slog.String("delivery_status", string(article.DeliveryStatus)))
	}

	if emailResp != nil && emailResp.EmailUUID != "" {
		addLogAttr(ctx, slog.String("email_uuid", emailResp.EmailUUID))
		addLogAttr(ctx, slog.String("email_provider", string(h.cfg.EmailProvider)))
	}

	return &message
}
