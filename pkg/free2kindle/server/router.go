package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/config"
)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// NewRouter creates and configures a new chi router with all middleware and routes.
func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	handlers := newHandlers(&handlerDeps{
		kindleEmail:      cfg.KindleEmail,
		senderEmail:      cfg.SenderEmail,
		mailjetAPIKey:    cfg.MailjetAPIKey,
		mailjetAPISecret: cfg.MailjetAPISecret,
	})

	r.Use(middleware.Recoverer)
	r.Use(requestIDMiddleware)
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)
	r.Use(jsonContentTypeMiddleware)

	r.Route("/api/v1", func(r chi.Router) {
		// Unauthenticated routes
		r.Get("/health", handlers.handleHealth)

		// Authenticated routes
		r.Route("/articles", func(r chi.Router) {
			r.Use(authMiddleware(cfg.APIKeySecret))
			r.Post("/", handlers.handleCreateArticle)
		})
	})

	return r
}
