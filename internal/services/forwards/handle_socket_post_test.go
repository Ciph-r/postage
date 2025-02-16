package forwards

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ciph-r/postage/internal/traffic"
	"github.com/stretchr/testify/require"
)

// TestHandleSocketPost checks that given a client is connected with an id of 1,
// when a BE service sends a post requests to it, then the request body is
// forwarded to the client, and the clients response is written back to the
// service.
func TestHandleSocketPost(t *testing.T) {
	loadBalancerMock := &LoadBalancerMock{
		SendSocketFunc: func(ctx context.Context, socketID string, r io.Reader) error {
			require.Equal(t, "1", socketID)
			require.Equal(t, "foo", mustReadStr(t, r))
			return nil
		},
	}
	mux := http.NewServeMux()
	HandleSocketPost(mux, loadBalancerMock)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/sockets/1", strings.NewReader("foo"))
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandleSocketPost_socket_id_not_open(t *testing.T) {
	loadBalancerMock := &LoadBalancerMock{
		SendSocketFunc: func(ctx context.Context, socketID string, r io.Reader) error {
			return traffic.ErrNotFound
		},
	}
	mux := http.NewServeMux()
	HandleSocketPost(mux, loadBalancerMock)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/sockets/1", strings.NewReader("foo"))
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func mustReadStr(t *testing.T, r io.Reader) string {
	t.Helper()
	b, err := io.ReadAll(r)
	require.NoError(t, err)
	return string(b)
}
