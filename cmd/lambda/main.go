// Lambda is the entry point for AWS Lambda deployment using API Gateway.
package main

import (
	"log/slog"
	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/shaftoe/savetoink/internal/config"
	"github.com/shaftoe/savetoink/internal/constant"
	"github.com/shaftoe/savetoink/internal/server"
)

func main() {
	cfg, err := config.Load(constant.ModeServer)
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	router := server.NewRouter(cfg)

	algnhsa.ListenAndServe(router, nil)
}
