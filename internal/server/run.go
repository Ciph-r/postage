package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v11"
	"github.com/ciph-r/postage/internal/services"
	"github.com/ciph-r/postage/internal/services/health"
	"github.com/ciph-r/postage/internal/services/sockets"
	"github.com/joho/godotenv"
)

func Run(ctx context.Context) error {
	slog.Info("server is starting")
	// gracefully handle ctrl-c by canceling the context.
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()
	// load config
	_ = godotenv.Load()
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	// build server dependencies.
	socketSvc := sockets.NewService(cfg.Sockets)
	healthSvc := health.NewService(cfg.Health)
	// run all the services.
	if err := services.RunGroup(ctx,
		socketSvc,
		healthSvc,
	); err != nil {
		return fmt.Errorf("failed to run services: %w", err)
	}
	slog.Info("server has gracefully stopped")
	return nil
}
