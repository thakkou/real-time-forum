package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func HandleClient(conn *websocket.Conn) {
	client := &Client{
		ID:   conn.RemoteAddr().String(),
		Conn: conn,
	}

	AddClient(client)
	fmt.Println("connected:", client.ID)

	go handle(client)
}

func handle(c *Client) {
	defer func() {
		RemoveClient(c.ID)
		c.Conn.Close()
		fmt.Println("disconnected:", c.ID)
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		processEvent(c, msg)
	}
}

func processEvent(c *Client, msg []byte) {
	switch string(msg) {

	case "ping":
		c.Conn.WriteMessage(websocket.TextMessage, []byte("pong"))

	case "hello":
		c.Conn.WriteMessage(websocket.TextMessage, []byte("hi"))

	default:
		Broadcast(c.ID, msg)
	}
}
