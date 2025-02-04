package services

import (
	"context"
	"time"
)

// httpServer is the interface of *http.Server and is here for testing
// purposes.
type httpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// NewHTTP turns any httpService into a Service. shutdownTimeout is the max
// amount of time the shutdown process will wait to gracefully shutdown.
func NewHTTP(srv httpServer, shutdownTimeout time.Duration) Service {
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
