package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ciph-r/postage/internal/server"
)

func main() {
	if err := server.Run(context.Background()); err != nil {
		slog.Error("server failed to run", "reason", err)
		os.Exit(1)
	}
}

func init() {
	// ensure all logs are json.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}
