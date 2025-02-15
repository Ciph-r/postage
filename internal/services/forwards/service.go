package forwards

import (
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
	"github.com/ciph-r/postage/internal/traffic"
)

func NewService(cfg Config, lb traffic.LoadBalancer) services.Service {
	mux := http.NewServeMux()
	HandleSocketPost(mux, lb)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
	return services.NewHTTP(srv, time.Minute)
}
