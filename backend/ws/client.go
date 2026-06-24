package ws

import "fmt"

func HandleClient(client *Client) {
	defer func() {
		// here if conction close i will send a event to front (disconect)
		mu.Lock()
		delete(Clients, client.id)
		mu.Unlock()

		client.conn.Close()
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
