package sockets

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func newHandleSocketServer(t *testing.T) * httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	HandleSocket(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func converToWSURL(target string) string {
	ws := strings.Replace(target, "http://", "ws://", 1)
	wss := strings.Replace(ws, "https://", "wss://", 1)
	return wss
}

func closeWS(t *testing.T, conn *websocket.Conn) {
	t.Helper()
	closeHandler := conn.CloseHandler()
	err := closeHandler(websocket.CloseNormalClosure, "")
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestHandleSocket(t *testing.T) {
	srv := newHandleSocketServer(t)
	baseUrl := converToWSURL(srv.URL)
	target, err := url.JoinPath(baseUrl, "/ws/testID")
	require.NoError(t, err)
	conn, _, err := websocket.DefaultDialer.Dial(target, nil)
	t.Cleanup(func() { closeWS(t, conn)})
	require.NoError(t, err)
}