package traffic

import (
	"context"
	"io"
)

// LoadBalancer intercepts requests for connected client.
type LoadBalancer interface {
	// OpenSocket binds a Socket to a given id. If the socket id is already
	// bound then ErrAlreadyExists is returned.
	OpenSocket(id string) (Socket, error)
	// SendSocket sends a reader to an open socket id, it if exists. if the
	// socket is not open it returns ErrNotFound.
	SendSocket(ctx context.Context, socketID string, r io.ReadCloser) error
}

type Socket interface {
	// Recv a io.ReadCloser. caller should loop this until done and then call
	// close.
	Recv() <-chan io.ReadCloser
	// defer this to clean up the socket in the load balancer.
	Close()
}
