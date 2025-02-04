package health

import (
	"log/slog"
	"net/http"
)

func Check(w http.ResponseWriter, r *http.Request) {
	slog.Debug("health check performed")
}
