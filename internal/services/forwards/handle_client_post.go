package forwards

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ciph-r/postage/internal/traffic"
)

// HandleClientPost allows a BE service to post data to a connected client by its client id.
func HandleClientPost[C traffic.ClientConnection](mux *http.ServeMux, clients traffic.ClientLoadBalancer[C]) {
	mux.HandleFunc("POST /clients/{clientID}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID := r.PathValue("clientID")
		// get connected client
		clientConn, err := clients.GetClientConnection(ctx, clientID)
		switch {
		case errors.Is(err, traffic.ErrNotFound):
			http.Error(w, "", http.StatusNotFound)
			return
		case err != nil:
			slog.Error("failed to get client connection", "clientID", clientID, "reason", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		// send data to client
		resp, err := clientConn.Post(r.Body)
		switch {
		case errors.Is(err, traffic.ErrDisconnected):
			http.Error(w, "", http.StatusNotFound)
			return
		case err != nil:
			slog.Warn("failed to post request body", "reason", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer resp.Close()
		// forward response back to caller
		switch n, err := io.Copy(w, resp); {
		case errors.Is(err, traffic.ErrDisconnected):
			http.Error(w, fmt.Sprintf("transfered %d bytes", n), http.StatusNotFound)
			return
		case err != nil:
			slog.Warn("failed to copy response body", "bytes written", n, "reason", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if err := resp.Close(); err != nil {
			slog.Warn("failed to close response body", "reason", err)
		}
	})
}
