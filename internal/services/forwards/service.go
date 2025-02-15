package forwards

import (
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
	"github.com/ciph-r/postage/internal/traffic"
)

func NewService(cfg Config, cc traffic.LoadBalancer) services.Service {
	mux := http.NewServeMux()
	HandleSocketPost(mux, cc)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
	return services.NewHTTP(srv, time.Minute)
}
