package ws

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type string          `json:"event_type"`
	Data json.RawMessage `json:"data"`
}
type TypingData struct {
	ConversationID int    `json:"conversationId"`
	ReceiverID     int    `json:"receiverId"`
	UserId         string `json:"userId"`
}

type Client struct {
	conn   *websocket.Conn
	isAuth bool
	id     string
}

var (
	Clients = make(map[string]map[*Client]bool)
	mu      sync.RWMutex
)

// this func send the notification and the data to all users exept u
func BroadcastExcept(senderID string, eventType string, data any) {
	payload := map[string]any{
		"event_type": eventType,
		"data":       data,
	}

	mu.RLock()
	defer mu.RUnlock()

	for userID, clients := range Clients {
		if userID == senderID {
			continue
		}

		for client := range clients {
			client.conn.WriteJSON(payload)
		}
	}
}

// this function send the notification to a special user
func NotifyUser(userID string, eventType string, data any) {
	mu.RLock()
	clients := Clients[userID]
	mu.RUnlock()

	payload := map[string]any{
		"event_type": eventType,
		"data":       data,
	}

	for client := range clients {
		client.conn.WriteJSON(payload)
	}
}

func RemoveClient(userID string, client *Client) {
	mu.Lock()
	defer mu.Unlock()

	if Clients[userID] != nil {
		delete(Clients[userID], client)

		if len(Clients[userID]) == 0 {
			delete(Clients, userID)
		}
	}
}

func StoreClient(userID string, conn *websocket.Conn) *Client {
	client := &Client{
		conn: conn,
		id:   userID,
	}

	mu.Lock()

	if Clients[userID] == nil {
		Clients[userID] = make(map[*Client]bool)
	}

	Clients[userID][client] = true

	// build online users list
	online := make([]string, 0)
	for id := range Clients {
		online = append(online, id)
	}

	mu.Unlock()

	NotifyUser(userID, "init", online)
	BroadcastExcept(userID, "client_connect", userID)

	return client
}

func HandleMessage(client *Client, raw []byte) {
	fmt.Println(string(raw))

	var msg WSMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Type: %s\n", msg.Type)
	fmt.Printf("Data: %s\n", string(msg.Data))

	switch msg.Type {
	case "new_posts": // for all users exepts u
		fmt.Println("new posts_notification")
	case "like_posts": // for u
		fmt.Println("user a liked ur posts")
	case "new_comments": // for u
		fmt.Println("new comments_notification")
	case "like_commnets": // for u
		fmt.Println("user a liked ur comments")
	case "send_message": // for u
		fmt.Println("message sent to user a")
	case "typing:start", "typing:stop":
		var data TypingData

		if err := json.Unmarshal(msg.Data, &data); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("data", data)

		NotifyUser(strconv.Itoa(data.ReceiverID), msg.Type, data)
	default:
		fmt.Println("unknown event:", msg.Type)
	}
}
