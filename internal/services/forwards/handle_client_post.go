package forwards

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/ciph-r/postage/internal/traffic"
)

func HandleClientPost[C traffic.ClientConnection](mux *http.ServeMux, clients traffic.ClientLoadBalancer[C]) {
	mux.HandleFunc("POST /clients/{clientID}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID := r.PathValue("clientID")
		clientConn, err := clients.GetClientConnection(ctx, clientID)
		switch {
		case errors.Is(err, traffic.ErrNotFound):
			http.Error(w, "specified client id is not connected", http.StatusNotFound)
			return
		case err != nil:
			slog.Error("failed to get client sender", "clientID", clientID, "reason", err)
			http.Error(w, "failed to get specified client", http.StatusInternalServerError)
			return
		}
		resp, err := clientConn.Post(r.Body)
		if err != nil {
			slog.Warn("failed to send request body", "reason", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer resp.Close()
		if n, err := io.Copy(w, resp); err != nil {
			slog.Warn("failed to copy response body", "bytes written", n, "reason", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if err := resp.Close(); err != nil {
			slog.Warn("failed to close response body", "reason", err)
		}
	})
}
