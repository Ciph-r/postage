package forwards

import (
	"net/http"

	"github.com/ciph-r/postage/internal/traffic"
)

// HandleClientPost allows a BE service to post data to a connected client by its client id.
func HandleClientPost(mux *http.ServeMux, clients traffic.LoadBalancer) {
	mux.HandleFunc("POST /clients/{clientID}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID := r.PathValue("clientID")
		// get connected client
		if err := clients.Forward(ctx, clientID, r.Body); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}
