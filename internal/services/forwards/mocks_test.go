package forwards

import (
	"context"
	"io"

	"github.com/ciph-r/postage/internal/traffic"
)

// ClientLoadBalancerMock satisfies traffic.ClientLoadBalancer interface.
type ClientLoadBalancerMock struct {
	// this is a hack to avoid mocking unused methods. it works similar to
	// method overloading in other languages. its allows the struct to satisfy
	// the interface without actually implementing the methods, or in this case,
	// only implementing the methods we need. this pattern is really only
	// acceptable in tests for ease of mocking.
	traffic.LoadBalancer

	ForwardFunc func(ctx context.Context, socketID string, r io.Reader) error
}

func (c *ClientLoadBalancerMock) Forward(ctx context.Context, socketID string, r io.Reader) error {
	return c.ForwardFunc(ctx, socketID, r)
}
