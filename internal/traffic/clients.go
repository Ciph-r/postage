package traffic

import (
	"context"
	"io"
)

// LoadBalancer intercepts requests for connected client.
type LoadBalancer interface {
	RegisterSocket(id string) (<-chan Forward, error)
	Forward(ctx context.Context, socketID string, r io.Reader) error
}

type Forward interface {
	io.Reader
}
