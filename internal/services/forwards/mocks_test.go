package forwards

import (
	"context"
	"io"

	"github.com/ciph-r/postage/internal/traffic"
)

// LoadBalancerMock satisfies traffic.ClientLoadBalancer interface.
type LoadBalancerMock struct {
	traffic.LoadBalancer
	SendSocketFunc func(ctx context.Context, socketID string, r io.ReadCloser) error
}

func (c *LoadBalancerMock) SendSocket(ctx context.Context, socketID string, r io.ReadCloser) error {
	return c.SendSocketFunc(ctx, socketID, r)
}
