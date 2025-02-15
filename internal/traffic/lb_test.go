package traffic

import (
	"context"
	"io"
)

func ExampleLoadBalancer_OpenSocket() {
	var lb LoadBalancer
	var ctx context.Context
	// open the socket
	socket, err := lb.OpenSocket("foo")
	if err != nil {
		return
	}
	defer socket.Close()
	// loop the socket
	for {
		var w io.Writer
		select {
		case <-ctx.Done():
			return
		case r := <-socket.Recv():
			_, _ = io.Copy(w, r)
		}
	}
}
