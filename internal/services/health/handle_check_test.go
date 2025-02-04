package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleCheck_healthy(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mux := http.NewServeMux()
	HandleCheck(mux)
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandleCheck_unhealthy(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mux := http.NewServeMux()
	checkers := []Checker{
		CheckFunc(func(_ context.Context) error {
			return errors.New("can't connect to db")
		}),
		CheckFunc(func(_ context.Context) error {
			return nil
		}),
	}
	HandleCheck(mux, checkers...)
	mux.ServeHTTP(w, r)
	require.NotEqual(t, http.StatusOK, w.Code)
}
