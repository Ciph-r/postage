package services

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// HTTPServerMock satisfies the httpServer interface.
type HTTPServerMock struct {
	ListenAndServeFunc func() error
	ShutdownFunc       func(ctx context.Context) error
}

func (h HTTPServerMock) ListenAndServe() error {
	return h.ListenAndServeFunc()
}
func (h HTTPServerMock) Shutdown(ctx context.Context) error {
	return h.ShutdownFunc(ctx)
}

// TestNewHTTP_stops_when_cancelled checks that when the Service context is
// cancelled the http servers Shutdown method is called to gracefully shutdown
// the server.
func TestNewHTTP_stops_when_cancelled(t *testing.T) {
	shutdownCalled := make(chan struct{})
	srv := HTTPServerMock{
		ListenAndServeFunc: func() error {
			<-shutdownCalled
			return http.ErrServerClosed
		},
		ShutdownFunc: func(ctx context.Context) error {
			close(shutdownCalled)
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	go cancel()
	err := NewHTTP(srv, time.Minute).Run(ctx)
	require.NoError(t, err)
}

// TestNewHTTP_stops_when_ListenAndServe_stops checks that the http Service
// will unblock correctly if ListenAndServe stops for any reason. the http
// servers Shutdown method should not be called in this case since the server is
// no longer running.
func TestNewHTTP_stops_when_ListenAndServe_stops(t *testing.T) {
	testErr := errors.New("boom")
	srv := HTTPServerMock{
		ListenAndServeFunc: func() error {
			return testErr
		},
	}
	err := NewHTTP(srv, time.Minute).Run(context.Background())
	require.ErrorIs(t, err, testErr)
}
