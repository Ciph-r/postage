package sockets

import (
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
)

func NewService(cfg Config) services.Service {
	mux := http.NewServeMux()
	HandleSocket(mux)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
	return services.NewHTTP(srv, time.Minute)
}
