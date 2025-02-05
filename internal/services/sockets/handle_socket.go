package sockets

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Globals
var clients = make(map[*websocket.Conn]bool)
var streamChan = make(chan []byte, 100)
var lock sync.Mutex
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //Allows all origins for now
	},
}

// HandleSocket will accept websocket connections on '/ws'
func HandleSocket(mux *http.ServeMux) {
	//Streams data to all connected clients
	//TODO: Make it possible to only broadcast to specific groups of users.
	go func() {
		for data := range streamChan {
			for conn := range clients {
				err := conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					slog.Error("Error sending to client", "reason", err)
					lock.Lock()
					delete(clients, conn)
					lock.Unlock()
				}
			}
		}
	}()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("Failed to upgrade connection", "reason", err)
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		lock.Lock()
		clients[conn] = true
		lock.Unlock()
		wg.Add(1)
		//Remove client when they disconnect
		go func() {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					lock.Lock()
					delete(clients, conn)
					lock.Unlock()
					wg.Done()
					break
				}
			}
		}()
		wg.Wait()
	})
}
