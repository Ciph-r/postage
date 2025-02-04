package health

import (
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
)

func NewService(cfg Config) services.Service {
	mux := http.NewServeMux()
	HandleCheck(mux)
	return services.NewHTTP(&http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}, time.Second)
}
