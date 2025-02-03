package server

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Service is a blocking process. It must gracefully shutdown if its context is
// cqnceled.
type Service interface {
	Run(ctx context.Context) error
}

// ServiceFunc satisfies the Service interface.
type ServiceFunc func(ctx context.Context) error

func (fn ServiceFunc) Run(ctx context.Context) error {
	return fn(ctx)
}

// runServices runs all specified services as a group under a given context. if
// one service is stopped, all other services are signaled to stop as well.
func runServices(ctx context.Context, svcs ...Service) error {
	errs := make([]error, len(svcs))
	var wg sync.WaitGroup
	signalStopped := make(chan struct{})
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if len(svcs) > 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			<-signalStopped
		}()
	}
	for i, svc := range svcs {
		wg.Add(1)
		go func(ctx context.Context, i int) {
			defer wg.Done()
			defer func() {
				select {
				case signalStopped <- struct{}{}:
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

// httpServer is the interface of *http.Server and is here for testing
// purposes.
type httpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// httpService turns any httpService into a Service. shutdownTimeout is the max
// amount of time the shutdown process will wait to gracefully shutdown.
func httpService(srv httpServer, shutdownTimeout time.Duration) Service {
	return ServiceFunc(func(ctx context.Context) error {
		srvStoppedErr := make(chan error)
		go func() {
			err := srv.ListenAndServe()
			select {
			case srvStoppedErr <- err:
			default:
			}
		}()
		select {
		case err := <-srvStoppedErr:
			return err
		case <-ctx.Done():
			timeoutCtx, timeourCancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer timeourCancel()
			return srv.Shutdown(timeoutCtx)
		}
	})
}

// ErrServicePanic is the error retuned when a service paniced and was its stack
// trace was recovered.
var ErrServicePanic = errors.New("service panic")
