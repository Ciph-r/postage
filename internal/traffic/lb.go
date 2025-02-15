package traffic

import (
	"context"
	"io"
)

// LoadBalancer manages access to socket connections.
type LoadBalancer interface {
	// OpenSocket binds a Socket to a given id. If the socket id is already
	// bound then ErrAlreadyExists is returned.
	OpenSocket(id string) (Socket, error)
	// SendSocket sends a reader to an open socket. If the
	// socket is not open it returns ErrNotFound.
	SendSocket(ctx context.Context, socketID string, r io.ReadCloser) error
}

// Socket is a long lived connection that is waited on to receive incoming
// messages in the form of io.ReadCloser's.
type Socket interface {
	// Recv a io.ReadCloser. caller should loop this until done and then call
	// close.
	Recv() <-chan io.ReadCloser
	// defer this to clean up the socket in the load balancer.
	Close()
}
