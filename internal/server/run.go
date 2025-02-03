package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/ciph-r/postage/internal/services/health"
	"github.com/ciph-r/postage/internal/services/sockets"
)

func Run(ctx context.Context) error {
	slog.Info("server is starting")
	// cancel the context if ctrl-c is signalled.
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()
	// build server dependencies.
	socketSrv := sockets.NewServer()
	socketSvc := httpService(socketSrv, time.Minute)
	healthSrv := health.NewServer()
	healthSvc := httpService(healthSrv, time.Second)
	// run all the services.
	if err := runServices(ctx, socketSvc, healthSvc); err != nil {
		return fmt.Errorf("failed to run services: %w", err)
	}
	slog.Info("server has stopped")
	return nil
}
