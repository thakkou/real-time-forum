package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var (
	clients = make(map[string]*Client)
	mu      sync.Mutex
)

// add client
func AddClient(c *Client) {
	mu.Lock()
	clients[c.ID] = c
	mu.Unlock()
}

// remove client
func RemoveClient(id string) {
	mu.Lock()
	delete(clients, id)
	mu.Unlock()
}

// broadcast message to all clients
func Broadcast(senderID string, msg []byte) {
	mu.Lock()
	defer mu.Unlock()

	for id, c := range clients {
		if id == senderID {
			continue
		}

		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			c.Conn.Close()
			delete(clients, id)
		}
	}
}
