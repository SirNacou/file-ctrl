package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirnacou/file-ctrl/backend/internal/api"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server, err := api.NewServer()
	if err != nil {
		slog.Error("init failed", slog.Any("error", err))
		os.Exit(1)
	}

	if err := server.Run(ctx); err != nil {
		slog.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
