package traffic

import (
	"context"
	"io"
)

// ClientLoadBalancer intercepts requests for connected client.
type ClientLoadBalancer[C ClientConnection] interface {
	GetClientConnection(ctx context.Context, id string) (C, error)
}

// ClientConnection is a conenct client that can recieve messages and respond.
type ClientConnection interface {
	Post(send io.Reader) (recv io.ReadCloser, err error)
}
