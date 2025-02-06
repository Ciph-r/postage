package sockets

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// newHandleEchoServer builds a httptest server that serves HandleEcho. it
// cleans up teh server when the test is finished.
func newHandleEchoServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	HandleEcho(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func convertToWSURL(target string) string {
	ws := strings.Replace(target, "http://", "ws://", 1)
	wss := strings.Replace(ws, "https://", "wss://", 1)
	return wss
}

// TestHandleEcho tests that HandleEcho returns whatever message is sent to it.
func TestHandleEcho(t *testing.T) {
	srv := newHandleEchoServer(t)
	// build url
	baseURL := convertToWSURL(srv.URL)
	target, err := url.JoinPath(baseURL, "/ws/echo")
	require.NoError(t, err)
	// connect to echo handler over the test server
	conn, _, err := websocket.DefaultDialer.Dial(target, nil)
	require.NoError(t, err)
	// send a message
	w, err := conn.NextWriter(websocket.TextMessage)
	require.NoError(t, err)
	_, err = fmt.Fprint(w, "hello")
	require.NoError(t, err)
	require.NoError(t, w.Close())
	// get the echo msg back
	typ, r, err := conn.NextReader()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, typ)
	// assert that the server echo'd the message back
	var buff bytes.Buffer
	_, err = io.Copy(&buff, r)
	require.NoError(t, err)
	require.Equal(t, "hello", buff.String())
}
