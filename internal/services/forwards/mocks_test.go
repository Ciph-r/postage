package forwards

import (
	"context"
	"io"

	"github.com/ciph-r/postage/internal/traffic"
)

// ClientLoadBalancerMock satisfies traffic.ClientLoadBalancer interface.
type ClientLoadBalancerMock struct {
	RegisterSocketFunc func(id string) (<-chan traffic.Forward, error)
	ForwardFunc        func(ctx context.Context, socketID string, r io.Reader) error
}

func (c *ClientLoadBalancerMock) RegisterSocket(id string) (<-chan traffic.Forward, error) {
	return c.RegisterSocketFunc(id)
}

func (c *ClientLoadBalancerMock) Forward(ctx context.Context, socketID string, r io.Reader) error {
	return c.ForwardFunc(ctx, socketID, r)
}
