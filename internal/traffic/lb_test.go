package traffic

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
)

func ExampleLoadBalancer_OpenSocket() {
	// mock LoadBalancer impl
	var lb LoadBalancer = LoadBalancerMock{
		OpenSocketFunc: func(id string) (Socket, error) {
			return SocketMock{
				RecvFunc: func() <-chan io.ReadCloser {
					// buffer 1 reader on the channel then close it to break the
					// for loop.
					out := make(chan io.ReadCloser, 1)
					out <- io.NopCloser(strings.NewReader("foo"))
					close(out)
					return out
				},
				CloseFunc: func() {},
			}, nil
		},
	}
	ctx := context.Background()
	// open the socket
	socket, err := lb.OpenSocket("1")
	if err != nil {
		return
	}
	defer socket.Close()
	// loop the socket
	var w bytes.Buffer
	defer fmt.Println(w.String())
	for {
		select {
		case <-ctx.Done():
			return
		case r, ok := <-socket.Recv():
			if !ok {
				return // chan is empty and closed.
			}
			defer r.Close()
			if _, err := io.Copy(&w, r); err != nil {
				return
			}
		}
	}
	// Outputs:
	// foo
}

type LoadBalancerMock struct {
	SendSocketFunc func(ctx context.Context, socketID string, r io.Reader) error
	OpenSocketFunc func(id string) (Socket, error)
}

func (c LoadBalancerMock) SendSocket(ctx context.Context, socketID string, r io.Reader) error {
	return c.SendSocketFunc(ctx, socketID, r)
}

func (c LoadBalancerMock) OpenSocket(id string) (Socket, error) {
	return c.OpenSocketFunc(id)
}

type SocketMock struct {
	RecvFunc  func() <-chan io.ReadCloser
	CloseFunc func()
}

func (s SocketMock) Recv() <-chan io.ReadCloser {
	return s.RecvFunc()
}
func (s SocketMock) Close() {
	s.CloseFunc()
}
