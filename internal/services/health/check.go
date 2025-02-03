package health

import (
	"log/slog"
	"net/http"
)

func NewServer(cfg Config) *http.Server {
	return &http.Server{
		Addr:    cfg.Addr,
		Handler: http.HandlerFunc(Check),
	}
}

func Check(w http.ResponseWriter, r *http.Request) {
	slog.Debug("health check performed")
}
