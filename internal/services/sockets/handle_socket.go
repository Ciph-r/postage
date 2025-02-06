package sockets

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Globals
var connectedClients = NewConnectedClients()
var streamChan = make(chan []byte, 100)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //Allows all origins for now
	},
}

// HandleSocket will accept websocket connections on '/ws/{clientID}'
func HandleSocket(mux *http.ServeMux) {
	mux.HandleFunc("GET /ws/{clientID}", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.PathValue("clientID")
		var wg sync.WaitGroup
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("Failed to upgrade connection", "reason", err)
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}

		if _, exists := connectedClients.GetClient(clientID); exists {
			slog.Error("Client already connected", "reason", err)
			return
		}
		connectedClients.AddClient(clientID, conn)
		defer connectedClients.DeleteClient(clientID)

		wg.Add(1)
		//Remove client when they disconnect
		go func() {
			for {
				_, _, err := connectedClients.conns[clientID].ReadMessage()
				if err != nil {
					connectedClients.DeleteClient(clientID)
					wg.Done()
					break
				}
			}
		}()
		//Send messages to connected client
		go func() {
			for {
				if conn, exists := connectedClients.GetClient(clientID); exists {
					err := conn.WriteMessage(websocket.TextMessage, <-streamChan)
					if err != nil {
						slog.Error("Error sending to client", "reason", err)
					}
				}
			}
		}()
		wg.Wait()
	})
}
