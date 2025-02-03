package health

import (
	"log/slog"
	"net/http"
)

func NewServer() *http.Server {
	return &http.Server{
		Addr:    ":9000",
		Handler: http.HandlerFunc(Check),
	}
}

func Check(w http.ResponseWriter, r *http.Request) {
	slog.Debug("health check performed")
}
