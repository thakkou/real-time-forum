package ws

import (
	"fmt"
)

func HandleClient(client *Client) {
	defer func() {
		// 1. Use your existing RemoveClient function instead of deleting the whole user map
		RemoveClient(client.id, client)
		client.conn.Close()

		// 2. Only broadcast disconnect if the user has no more active tabs open
		mu.RLock()
		_, stillOnline := Clients[client.id]
		mu.RUnlock()

		if !stillOnline {
			BroadcastExcept(client.id, "client_disconnect", client.id)
		}
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Println("client disconnected:", client.id)
			return
		}

		HandleMessage(client, msg)
	}
}
