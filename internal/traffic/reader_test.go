package traffic

import (
	"fmt"
	"strings"
)

func Example_waitingReadCloser() {
	r := newWaitingReadCloser(strings.NewReader("foo"))
	go func() {
		defer r.Close()
		fmt.Fscanln(r)
	}()
	// wait for Fscanln to finish.
	<-r.Wait()
	// Outputs: foo
}
