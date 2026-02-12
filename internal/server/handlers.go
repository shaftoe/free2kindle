// Package server provides HTTP handlers and middleware for the savetoink application.
package server

import (
	"github.com/shaftoe/savetoink/internal/config"
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
