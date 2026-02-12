package server

import (
	"log/slog"

	"github.com/shaftoe/savetoink/internal/config"
	"github.com/shaftoe/savetoink/internal/model"
	"github.com/shaftoe/savetoink/internal/service"
)

type articleRequest struct {
	URL string `json:"url"`
}

type articleResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	Message        string `json:"message"`
	DeliveryStatus string `json:"delivery_status,omitempty"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type listArticlesResponse struct {
	Articles []*model.Article `json:"articles"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	Total    int              `json:"total"`
	HasMore  bool             `json:"has_more"`
}

type deleteArticleResponse struct {
	Deleted int `json:"deleted"`
}

type handlers struct {
	cfg     *config.Config
	service service.Interface
}

type contextKey string

type logRecord struct {
	*slog.Record
}
