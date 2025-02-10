package forwards

import (
	"context"
	"io"
)

// ClientLoadBalancerMock satisfies traffic.ClientLoadBalancer interface.
type ClientLoadBalancerMock struct {
	GetClientConnectionFunc func(ctx context.Context, id string) (*ClientConnectionMock, error)
}

func (c *ClientLoadBalancerMock) GetClientConnection(ctx context.Context, id string) (*ClientConnectionMock, error) {
	return c.GetClientConnectionFunc(ctx, id)
}

// ClientConnectionMock satisfies the clientConn interface.
type ClientConnectionMock struct {
	PostFunc func(send io.Reader) (recv io.ReadCloser, err error)
}

func (c *ClientConnectionMock) Post(send io.Reader) (recv io.ReadCloser, err error) {
	return c.PostFunc(send)
}
