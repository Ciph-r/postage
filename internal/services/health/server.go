package health

import "net/http"

func NewServer(cfg Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Check)
	return &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
}
