package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// ServiceMock satisfies the Service interface.
type ServiceMock struct {
	RunFunc func(ctx context.Context) error
}

func (s *ServiceMock) Run(ctx context.Context) error {
	return s.RunFunc(ctx)
}

// TestRunGroup_none test that the calling RunGroup without any services
// is no-op, and returns no error.
func TestRunGroup_none(t *testing.T) {
	require.NoError(t, RunGroup(context.Background()))
}

// TestRunGroup_respects_context tests that the context is passed to the
// service and that it respects when it is canceled.
func TestRunGroup_respects_context(t *testing.T) {
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
	err := RunGroup(ctx, svcA, svcB)
	require.ErrorIs(t, err, context.Canceled)
}

// TestRunGroup_stop_all_for_one_err checks when a single service stops,
// than all other services are signalled to stop as well.
func TestRunGroup_stop_all_for_one_err(t *testing.T) {
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
	err := RunGroup(context.Background(), svcA, svcB, svcC)
	require.Error(t, err)
	require.ErrorIs(t, err, testErr)
	require.ErrorIs(t, err, context.Canceled)
}

// TestRunGroup_stop_all_for_one_panic checks when a single service panics,
// the panic is recovered from and converted into an error, and than all other
// services are signalled to stop as well.
func TestRunGroup_stop_all_for_one_panic(t *testing.T) {
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
	err := RunGroup(context.Background(), svcA, svcB, svcC)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrServicePanic)
	require.ErrorIs(t, err, context.Canceled)
}
