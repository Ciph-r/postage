package health

import (
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
)

func NewService(cfg Config, checkers ...Checker) services.Service {
	mux := http.NewServeMux()
	HandleCheck(mux, checkers...)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
	return services.NewHTTP(srv, time.Second)
}
