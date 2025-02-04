package health

import (
	"log/slog"
	"net/http"
)

func HandleCheck(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("health check performed")
	})
}
