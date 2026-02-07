// HTTP server is the entry point for running the application as a standalone HTTP server.
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/server"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 15 * time.Second
)

func main() {
	cfg, err := config.Load(config.ModeServer)
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	if cfg.Debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	router := server.NewRouter(cfg)

	port := "8080"
	slog.Info("starting HTTP server", "port", port)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
	if srvErr := srv.ListenAndServe(); srvErr != nil {
		slog.Error("failed to start server", "error", srvErr)
		os.Exit(1)
	}
}
