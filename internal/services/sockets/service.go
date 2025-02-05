package sockets

import (
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
)

func NewService(cfg Config) services.Service {

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: nil, //TODO: set handler with socket router.
	}
	return services.NewHTTP(srv, time.Minute)
}
