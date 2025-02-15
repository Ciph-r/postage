package traffic

import (
	"io"
)

// ClientLoadBalancer intercepts requests for connected client.
type ClientLoadBalancer interface {
	RegisterSocket(id string, conn websocketConn)
	GetSocketReader() <-chan io.ReadCloser

	RegisterForward(is string, r io.ReadCloser)
	GetForwardReader() <-chan io.ReadCloser
}

type websocketConn interface {
	NextReader() (messageType int, r io.Reader, err error)
	NextWriter(messageType int) (io.WriteCloser, error)
}
