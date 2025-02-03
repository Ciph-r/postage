package server

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ServiceMock satisfies the Service interface.
type ServiceMock struct {
	RunFunc func(ctx context.Context) error
}

func (s *ServiceMock) Run(ctx context.Context) error {
	return s.RunFunc(ctx)
}

// Test_runServices_none test that the calling runServices without any services
// is no-op, and returns no error.
func Test_runServices_none(t *testing.T) {
	require.NoError(t, runServices(context.Background()))
}

// Test_runServices_respects_context tests that the context is passed to the
// service and that it respects when it is canceled.
func Test_runServices_respects_context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	bIsStarted := make(chan struct{})
	svcA := &ServiceMock{
		RunFunc: func(ctx context.Context) error {
			<-bIsStarted
			go cancel()
			<-ctx.Done()
			return ctx.Err()
		},
	}
	svcB := &ServiceMock{
		RunFunc: func(ctx context.Context) error {
			close(bIsStarted)
			<-ctx.Done()
			return ctx.Err()
		},
	}
	err := runServices(ctx, svcA, svcB)
	require.ErrorIs(t, err, context.Canceled)
}

// Test_runServices_stop_all_for_one_err checks when a single service stops,
// than all other services are signalled to stop as well.
func Test_runServices_stop_all_for_one_err(t *testing.T) {
	testErr := errors.New("boom")
	bIsStarted := make(chan struct{})
	cIsStarted := make(chan struct{})
	svcA := &ServiceMock{
		RunFunc: func(ctx context.Context) error {
			<-cIsStarted
			return testErr
		},
	}
	svcB := &ServiceMock{
		RunFunc: func(ctx context.Context) error {
			close(bIsStarted)
			<-ctx.Done()
			return ctx.Err()
		},
	}
	svcC := &ServiceMock{
		RunFunc: func(ctx context.Context) error {
			<-bIsStarted
			close(cIsStarted)
			<-ctx.Done()
			return ctx.Err()
		},
	}
	err := runServices(context.Background(), svcA, svcB, svcC)
	require.Error(t, err)
	require.ErrorIs(t, err, testErr)
	require.ErrorIs(t, err, context.Canceled)
}

// Test_runServices_stop_all_for_one_panic checks when a single service panics,
// the panic is recovered from and converted into an error, and than all other
// services are signalled to stop as well.
func Test_runServices_stop_all_for_one_panic(t *testing.T) {
	bIsStarted := make(chan struct{})
	cIsStarted := make(chan struct{})
	svcA := &ServiceMock{
		RunFunc: func(ctx context.Context) error {
			<-cIsStarted
			panic("boom")
		},
	}
	svcB := &ServiceMock{
		RunFunc: func(ctx context.Context) error {
			close(bIsStarted)
			<-ctx.Done()
			return ctx.Err()
		},
	}
	svcC := &ServiceMock{
		RunFunc: func(ctx context.Context) error {
			<-bIsStarted
			close(cIsStarted)
			<-ctx.Done()
			return ctx.Err()
		},
	}
	err := runServices(context.Background(), svcA, svcB, svcC)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrServicePanic)
	require.ErrorIs(t, err, context.Canceled)
}

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

func Test_httpService_stops_when_canceled(t *testing.T) {
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
	err := httpService(srv, time.Minute).Run(ctx)
	require.NoError(t, err)
}

func Test_httpService_unblocks_when_shutdown(t *testing.T) {
	testErr := errors.New("boom")
	srv := HTTPServerMock{
		ListenAndServeFunc: func() error {
			return testErr
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	go cancel()
	err := httpService(srv, time.Minute).Run(ctx)
	require.ErrorIs(t, err, testErr)
}
