package forwards

import (
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
	"github.com/ciph-r/postage/internal/traffic"
)

func NewService[C traffic.ClientConnection](cfg Config, cc traffic.ClientLoadBalancer[C]) services.Service {
	mux := http.NewServeMux()
	HandleClientPost(mux, cc)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
	return services.NewHTTP(srv, time.Minute)
}
