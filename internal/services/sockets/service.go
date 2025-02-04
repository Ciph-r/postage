package sockets

import (
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
)

func NewService() services.Service {
	srv := &http.Server{
		Addr:    ":80", // TODO: set addr with a configurable var.
		Handler: nil,   // TODO: set handler with socket router.
	}
	return services.NewHTTP(srv, time.Minute)
}
