package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrServicePanic is the error returned when a service panics and was its stack
// trace was recovered.
var ErrServicePanic = errors.New("service panic")

// Service is a blocking process. It must gracefully shutdown if its context is
// canceled.
type Service interface {
	Run(ctx context.Context) error
}

// ServiceFunc satisfies the Service interface.
type ServiceFunc func(ctx context.Context) error

func (fn ServiceFunc) Run(ctx context.Context) error {
	return fn(ctx)
}

// RunGroup runs all specified services as a group under a given context. if
// one service is stopped, all other services are signaled to stop as well.
func RunGroup(ctx context.Context, svcs ...Service) error {
	if len(svcs) < 1 {
		return ctx.Err()
	}
	errs := make([]error, len(svcs))
	var wg sync.WaitGroup
	oneServiceHasStopped := make(chan struct{})
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		<-oneServiceHasStopped
	}()
	for i, svc := range svcs {
		wg.Add(1)
		go func(ctx context.Context, i int) {
			defer wg.Done()
			defer func() {
				select {
				case oneServiceHasStopped <- struct{}{}:
				default:
				}
			}()
			defer func() {
				if v := recover(); v != nil {
					errs[i] = errors.Join(ErrServicePanic, fmt.Errorf("%v", v))
				}
			}()
			errs[i] = svc.Run(ctx)
		}(ctx, i)
	}
	wg.Wait()
	return errors.Join(errs...)
}
