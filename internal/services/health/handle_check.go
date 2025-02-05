package health

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"
)

// HandleCheck is used to inform deployment services of the health of the Postage server.
func HandleCheck(mux *http.ServeMux, checkers ...Checker) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("health check performed")
		ctx := r.Context()
		errs := make([]error, len(checkers))
		var wg sync.WaitGroup
		for i, c := range checkers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				errs[i] = c.Check(ctx)
			}()
		}
		wg.Wait()
		w.Header().Set("Content-Type", "application/json")
		if err := errors.Join(errs...); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(CheckResponse{
				Status:  Unhealthy,
				Message: err.Error(),
			})
		}
		json.NewEncoder(w).Encode(CheckResponse{Status: Healthy})
	})
}

type CheckResponse struct {
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
}
