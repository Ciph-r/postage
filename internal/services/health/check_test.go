package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleCheck(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mux := http.NewServeMux()
	HandleCheck(mux)
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
}
