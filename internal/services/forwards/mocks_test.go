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

// ClientConnectionMock satisfies the traffic.ClientConnection interface.
type ClientConnectionMock struct {
	PostFunc func(send io.Reader) (recv io.ReadCloser, err error)
}

func (c *ClientConnectionMock) Post(send io.Reader) (recv io.ReadCloser, err error) {
	return c.PostFunc(send)
}

// ReadCloserMock satisfies the ReadCloser interface.
type ReadCloserMock struct {
	CloseFunc func() error
	ReadFunc  func(p []byte) (n int, err error)
}

func (r *ReadCloserMock) Close() error {
	return r.CloseFunc()
}
func (r *ReadCloserMock) Read(p []byte) (n int, err error) {
	return r.ReadFunc(p)
}
