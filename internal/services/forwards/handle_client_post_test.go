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

func TestHandleClientPost(t *testing.T) {
	clientConnectionMock := &ClientConnectionMock{
		PostFunc: func(send io.Reader) (recv io.ReadCloser, err error) {
			return io.NopCloser(strings.NewReader("bar")), nil
		},
	}
	clientLoadBalancerMock := &ClientLoadBalancerMock{
		GetClientConnectionFunc: func(ctx context.Context, id string) (*ClientConnectionMock, error) {
			require.Equal(t, "1", id)
			return clientConnectionMock, nil
		},
	}
	mux := http.NewServeMux()
	HandleClientPost(mux, clientLoadBalancerMock)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/clients/1", strings.NewReader("foo"))
	mux.ServeHTTP(w, r)
	// check response
	b, err := io.ReadAll(w.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "bar", string(b))
}

func TestHandleClientPost_not_connected(t *testing.T) {
	clientLoadBalancerMock := &ClientLoadBalancerMock{
		GetClientConnectionFunc: func(ctx context.Context, id string) (*ClientConnectionMock, error) {
			return nil, traffic.ErrNotFound
		},
	}
	mux := http.NewServeMux()
	HandleClientPost(mux, clientLoadBalancerMock)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/clients/1", strings.NewReader("foo"))
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleClientPost_disconencted_during_send(t *testing.T) {
	clientConnectionMock := &ClientConnectionMock{
		PostFunc: func(send io.Reader) (recv io.ReadCloser, err error) {
			return nil, traffic.ErrDisconnected
		},
	}
	clientLoadBalancerMock := &ClientLoadBalancerMock{
		GetClientConnectionFunc: func(ctx context.Context, id string) (*ClientConnectionMock, error) {
			return clientConnectionMock, nil
		},
	}
	mux := http.NewServeMux()
	HandleClientPost(mux, clientLoadBalancerMock)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/clients/1", strings.NewReader("foo"))
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleClientPost_disconencted_during_recv(t *testing.T) {
	readCloserMock := &ReadCloserMock{
		ReadFunc: func(p []byte) (n int, err error) {
			return 0, traffic.ErrDisconnected
		},
		CloseFunc: func() error {
			return nil
		},
	}
	clientConnectionMock := &ClientConnectionMock{
		PostFunc: func(send io.Reader) (recv io.ReadCloser, err error) {
			return readCloserMock, nil
		},
	}
	clientLoadBalancerMock := &ClientLoadBalancerMock{
		GetClientConnectionFunc: func(ctx context.Context, id string) (*ClientConnectionMock, error) {
			return clientConnectionMock, nil
		},
	}
	mux := http.NewServeMux()
	HandleClientPost(mux, clientLoadBalancerMock)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/clients/1", strings.NewReader("foo"))
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}
