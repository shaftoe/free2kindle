package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/config"
	f2khttp "github.com/shaftoe/free2kindle/pkg/free2kindle/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	if cfg.Debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	router := f2khttp.NewRouter(cfg)

	port := "8080"
	slog.Info("starting HTTP server", "port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
