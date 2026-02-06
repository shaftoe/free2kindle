package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/config"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	handlers := newHandlers(&handlerDeps{
		KindleEmail:      cfg.KindleEmail,
		SenderEmail:      cfg.SenderEmail,
		MailjetAPIKey:    cfg.MailjetAPIKey,
		MailjetAPISecret: cfg.MailjetAPISecret,
	})

	r.Use(middleware.Recoverer)
	r.Use(requestIDMiddleware)
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)
	r.Use(jsonContentTypeMiddleware)

	r.Route("/api/v1", func(r chi.Router) {
		// Unauthenticated routes
		r.Get("/health", handlers.HandleHealth)

		// Authenticated routes
		r.Route("/articles", func(r chi.Router) {
			r.Use(authMiddleware(cfg.APIKeySecret))
			r.Post("/", handlers.HandleCreateArticle)
		})
	})

	return r
}
