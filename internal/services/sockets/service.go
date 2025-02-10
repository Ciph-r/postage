package sockets

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ciph-r/postage/internal/services"
	"github.com/gorilla/websocket"
)

func NewService(cfg Config) services.Service {
	mux := http.NewServeMux()
	HandleEcho(mux)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
	return services.NewHTTP(srv, time.Minute)
}

func HandleEcho(mux *http.ServeMux) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true //Allows all origins for now
		},
	}
	mux.HandleFunc("GET /ws/echo", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := echo(ctx, conn); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func echo(ctx context.Context, conn *websocket.Conn) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("exiting echo loop: %w", ctx.Err())
		default:
		}
		messageType, r, err := conn.NextReader()
		if err != nil {
			return fmt.Errorf("failed to get next reader: %w", err)
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			return fmt.Errorf("failed to get next writer: %w", err)
		}
		if _, err := io.Copy(w, r); err != nil {
			return fmt.Errorf("failed to get copy: %w", err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}
	}
}
