package http

import (
	"encoding/json"
	"net/http"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email/mailjet"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/service"
)

type handlerDeps struct {
	KindleEmail      string
	SenderEmail      string
	MailjetAPIKey    string
	MailjetAPISecret string
}

type handlers struct {
	deps *handlerDeps
}

func newHandlers(deps *handlerDeps) *handlers {
	return &handlers{deps: deps}
}

func (h *handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

func (h *handlers) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	var req ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid request body"})
		return
	}

	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "URL is required"})
		return
	}

	mailjetConfig := &mailjet.Config{
		APIKey:      h.deps.MailjetAPIKey,
		APISecret:   h.deps.MailjetAPISecret,
		SenderEmail: h.deps.SenderEmail,
	}

	svcCfg := &service.Config{
		Extractor:    content.NewExtractor(),
		Generator:    epub.NewGenerator(),
		Sender:       mailjet.NewSender(mailjetConfig),
		SendEmail:    true,
		GenerateEPUB: true,
		KindleEmail:  h.deps.KindleEmail,
		SenderEmail:  h.deps.SenderEmail,
		Subject:      "",
		OutputPath:   "",
	}

	result, err := service.Run(r.Context(), svcCfg, req.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Failed to process article: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ArticleResponse{
		Title:   result.Title,
		URL:     req.URL,
		Status:  "completed",
		Message: "article sent to Kindle successfully",
	})
}
