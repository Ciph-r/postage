package forwards

import (
	"context"
	"io"

	"github.com/ciph-r/postage/internal/traffic"
)

// ClientLoadBalancerMock satisfies traffic.ClientLoadBalancer interface.
type ClientLoadBalancerMock struct {
	traffic.LoadBalancer

	ForwardFunc func(ctx context.Context, socketID string, r io.Reader) error
}

func (c *ClientLoadBalancerMock) Forward(ctx context.Context, socketID string, r io.Reader) error {
	return c.ForwardFunc(ctx, socketID, r)
}
