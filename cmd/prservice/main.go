package main

import (
	"context"
	"log/slog"

	"github.com/Egorrrad/avitotechBackendPR/internal/app/prservice"
)

func main() {
	ctx := context.Background()
	cfg, err := prservice.LoadConfig()
	if err != nil {
		slog.Error("Load config", "error", err)
	}

	prservice.InitLogging()

	server, err := prservice.NewServer(cfg)
	if err != nil {
		slog.Error("APIServer init", "error", err)
	}

	server.Start(ctx)
}
