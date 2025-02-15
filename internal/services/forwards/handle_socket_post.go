package forwards

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ciph-r/postage/internal/traffic"
)

// HandleSocketPost allows a BE service to post data to a connected client by its client id.
func HandleSocketPost(mux *http.ServeMux, lb traffic.LoadBalancer) {
	mux.HandleFunc("POST /sockets/{socketID}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		socketID := r.PathValue("socketID")
		switch err := lb.SendSocket(ctx, socketID, r.Body); {
		case errors.Is(err, traffic.ErrNotFound):
			http.Error(w, "", http.StatusNotFound)
			return
		case err != nil:
			slog.Error("failed to send to socket", "socketID", socketID, "reason", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}
