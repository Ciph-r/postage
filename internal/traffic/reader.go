package traffic

import (
	"io"
	"sync/atomic"
)

// waitingReadCloser satisfies the io.ReadCloser interface.
type waitingReadCloser struct {
	r        io.Reader
	n        int
	done     chan struct{}
	isClosed atomic.Bool
}

func newWaitingReadCloser(r io.Reader) *waitingReadCloser {
	return &waitingReadCloser{
		r:    r,
		done: make(chan struct{}, 1),
	}
}

// Wait for Close to signal done, too allow resourse cleanup.
func (w *waitingReadCloser) Wait() <-chan struct{} {
	return w.done
}

func (w *waitingReadCloser) Read(p []byte) (n int, err error) {
	if w.isClosed.Load() {
		return w.n, io.EOF
	}
	w.n += len(p)
	return w.r.Read(p)
}

// Close can safely be called mutliple times.
func (w *waitingReadCloser) Close() error {
	w.isClosed.Store(true)
	select {
	case w.done <- struct{}{}:
	default:
	}
	return nil
}
