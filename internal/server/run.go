package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/ciph-r/postage/internal/services/sockets"
)

func Run(ctx context.Context) error {
	slog.Info("started")
	// cancel the context if ctrl-c is signalled.
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()
	// build server dependencies.
	socketSrv := sockets.NewServer()
	socketSvc := httpService(socketSrv, time.Minute)
	// run all the services.
	if err := runServices(ctx, socketSvc); err != nil {
		return fmt.Errorf("failed to run services: %w", err)
	}
	slog.Info("stopped")
	return nil
}
