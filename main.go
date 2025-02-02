package main

import (
	"log/slog"
	"os"

	"github.com/ciph-r/postage/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		slog.Error("server failed to run", "reason", err)
		os.Exit(1)
	}
}

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}
