package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shaftoe/free2kindle/internal/auth"
	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/constant"
	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/email"
	"github.com/shaftoe/free2kindle/internal/email/mailjet"
	"github.com/shaftoe/free2kindle/internal/epub"
	"github.com/shaftoe/free2kindle/internal/repository"
	"github.com/shaftoe/free2kindle/internal/service"
)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// NewRouter creates and configures a new chi router with all middleware and routes.
func NewRouter(cfg *config.Config) *chi.Mux {
	setupLogging(cfg)

	r := chi.NewRouter()
	srv := newService(cfg)
	handlers := newHandlers(
		cfg,
		srv,
		repository.NewDynamoDB(cfg.AWSConfig, cfg.DynamoDBTable),
	)

	r.Use(middleware.Recoverer)
	r.Use(auth.NewUserIDMiddleware(cfg))
	r.Use(requestIDMiddleware)
	r.Use(corsMiddleware)
	r.Use(jsonContentTypeMiddleware)
	r.Use(loggingMiddleware)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "not_found"})
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "method_not_allowed"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handlers.handleHealth)

		r.Route("/articles", func(r chi.Router) {
			r.Use(auth.EnsureAutheticatedMiddleware)
			r.Post("/", handlers.handleCreateArticle)
		})
	})

	return r
}

func setupLogging(cfg *config.Config) {
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}

func newService(cfg *config.Config) *service.Service {
	var sender email.Sender
	if cfg.SendEnabled {
		switch cfg.EmailProvider {
		case constant.EmailBackendMailjet:
			sender = mailjet.NewSender(cfg.MailjetAPIKey, cfg.MailjetAPISecret, cfg.SenderEmail)
		default:
			sender = mailjet.NewSender(cfg.MailjetAPIKey, cfg.MailjetAPISecret, cfg.SenderEmail)
		}
	}

	return service.New(service.NewDeps(
		content.NewExtractor(),
		epub.NewGenerator(),
		sender,
	), cfg)
}
