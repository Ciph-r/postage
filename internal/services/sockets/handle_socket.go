package sockets

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Globals
var connectedClients = NewConnectedClients()

// HandleSocket will accept websocket connections on '/ws/{clientID}'
func HandleSocket(mux *http.ServeMux) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true //Allows all origins for now
		},
	}
	mux.HandleFunc("GET /ws/{clientID}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clientID := r.PathValue("clientID")
		var wg sync.WaitGroup

		if _, exists := connectedClients.GetClient(clientID); exists {
			slog.Error("Client already connected", "reason", err)
			return
		}
		connectedClients.AddClient(clientID, conn)
		defer connectedClients.DeleteClient(clientID)

		wg.Add(1)
		//Remove client when they disconnect
		go func() {
			defer wg.Done()
			if conn, exists := connectedClients.GetClient(clientID); exists {
				for {
					select {
					case <-ctx.Done():
						connectedClients.DeleteClient(clientID)
						return
					default:
					}
					_, _, err := conn.ReadMessage()
					if err != nil {
						connectedClients.DeleteClient(clientID)
						break
					}
				}
			}
		}()
		wg.Wait()
	})
}

// Replaces the channel that was previously used. Sends a binary message to client with {clientID} and returns a ReadCloser
// TODO: This functionality needs to be adapted into the Post function used by the traffic client interface
func sendToClient(clientId string, r io.Reader) (io.ReadCloser, error) {
	if conn, exists := connectedClients.GetClient(clientId); exists {
		w, err := conn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(w, r)
		if err != nil {
			return nil, err
		}
		w.Close()
		_, reader, err := conn.NextReader()
		if err != nil {
			return nil, err
		}
		return io.NopCloser(reader), nil
	}
	return nil, errors.New("client not found")
}
