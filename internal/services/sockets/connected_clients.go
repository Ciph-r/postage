package sockets

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectedClients struct {
	lock  sync.Mutex
	conns map[string]*websocket.Conn
}

func NewConnectedClients() *ConnectedClients {
	return &ConnectedClients{
		conns: make(map[string]*websocket.Conn),
	}
}

func (c *ConnectedClients) AddClient(id string, conn *websocket.Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.conns[id] = conn
}

func (c *ConnectedClients) GetClient(id string) (*websocket.Conn, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	conn, exists := c.conns[id]
	return conn, exists
}

func (c *ConnectedClients) DeleteClient(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.conns[id]; exists {
		c.conns[id].Close()
		delete(c.conns, id)
	}
}
