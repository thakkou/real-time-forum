package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"forum/database"
	// "forum/ws"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type string          `json:"event_type"`
	Data json.RawMessage `json:"data"`
}

// duplicate for : SendMessageRequest
type MessageCreationData struct {
	ConversationID *int   `json:"conversationId"`
	ReceiverID     int    `json:"receiverId"`
	Text           string `json:"text"`
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
	Clients = make(map[string]*Client)
	mu      sync.RWMutex
)

// this func send the notification and the data to all users exept u
func BroadcastExcept(senderID string, eventType string, data any) {
	fmt.Println("start broadcasting")

	payload := map[string]any{
		"event_type": eventType,
		"data":       data,
	}

	mu.RLock()
	clientsCopy := make([]*Client, 0, len(Clients))

	for userID, client := range Clients {
		if userID == senderID {
			continue
		}
		clientsCopy = append(clientsCopy, client)
	}
	mu.RUnlock()

	for _, client := range clientsCopy {
		if err := client.conn.WriteJSON(payload); err != nil {
			fmt.Println("broadcast error:", err)
		}
	}
}

// this function send the notification to a special user
func NotifyUser(userID string, eventType string, data any) {
	mu.RLock()
	client, ok := Clients[userID]
	mu.RUnlock()

	if !ok {
		return
	}

	payload := map[string]any{
		"event_type": eventType,
		"data":       data,
	}

	client.conn.WriteJSON(payload)
}

func StoreClient(userID string, conn *websocket.Conn) *Client {
	client := &Client{
		conn: conn,
		id:   userID,
	}

	mu.Lock()
	Clients[userID] = client
	// get Online users
	online := make([]string, 0, len(Clients))
	for id := range Clients {
		online = append(online, id)
	}
	mu.Unlock()
	NotifyUser(client.id, "init", online)
	BroadcastExcept(client.id, "client_connect", client.id)

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
	case "create_message": // for u
		var data MessageCreationData

		if err := json.Unmarshal(msg.Data, &data); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("data", data)

		// handle create message with data like did in the route !
		senderId, _ := strconv.Atoi(client.id)
		handleMessageCreation(senderId, data)

		// NotifyUser(strconv.Itoa(data.ReceiverID), msg.Type, data)
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

func handleMessageCreation(senderID int, data MessageCreationData) {
	fmt.Println("========== CREATE MESSAGE START ==========")

	fmt.Printf("[AUTH] sender=%d\n", senderID)

	// -------------------------
	// Validate
	// -------------------------
	if data.ReceiverID == 0 || data.Text == "" {
		fmt.Println("[VALIDATION] missing fields")
		return
	}

	if senderID == data.ReceiverID {
		fmt.Println("[VALIDATION] user tried to message himself")
		return
	}

	// -------------------------
	// Normalize pair
	// -------------------------
	user1 := senderID
	user2 := data.ReceiverID

	if user1 > user2 {
		user1, user2 = user2, user1
	}

	fmt.Printf("[CONVERSATION] normalized pair=(%d,%d)\n", user1, user2)

	// -------------------------
	// Start transaction
	// -------------------------
	tx, err := database.Database.Begin()
	if err != nil {
		fmt.Println("[DB] begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	var conversationID int

	// -------------------------
	// CASE 1: conversation_id provided
	// -------------------------
	if data.ConversationID != nil {

		conversationID = *data.ConversationID

		fmt.Printf("[CONVERSATION] validating conversation_id=%d\n", conversationID)

		var exists int

		err := tx.QueryRow(`
			SELECT id
			FROM CONVERSATIONS
			WHERE id = ?
			AND user1_id = ?
			AND user2_id = ?
		`,
			conversationID,
			user1,
			user2,
		).Scan(&exists)
		if err != nil {
			fmt.Println("[CONVERSATION] invalid conversation:", err)
			return
		}

		fmt.Printf("[CONVERSATION] validated id=%d\n", exists)

	} else {

		// -------------------------
		// CASE 2: Find or create conversation
		// -------------------------
		fmt.Printf("[CONVERSATION] searching (%d,%d)\n", user1, user2)

		err := tx.QueryRow(`
			SELECT id
			FROM CONVERSATIONS
			WHERE user1_id = ?
			AND user2_id = ?
		`,
			user1,
			user2,
		).Scan(&conversationID)

		if err == sql.ErrNoRows {

			fmt.Printf("[CONVERSATION] not found, creating (%d,%d)\n", user1, user2)

			res, err := tx.Exec(`
				INSERT INTO CONVERSATIONS (
					user1_id,
					user2_id
				)
				VALUES (?, ?)
			`,
				user1,
				user2,
			)
			if err != nil {
				fmt.Println("[CONVERSATION] create error:", err)
				return
			}

			id, err := res.LastInsertId()
			if err != nil {
				fmt.Println("[CONVERSATION] last insert id error:", err)
				return
			}

			conversationID = int(id)

			fmt.Printf("[CONVERSATION] created id=%d\n", conversationID)

		} else if err != nil {
			fmt.Println("[CONVERSATION] lookup error:", err)
			return
		} else {
			fmt.Printf("[CONVERSATION] found id=%d\n", conversationID)
		}
	}

	// -------------------------
	// Insert message
	// -------------------------
	fmt.Printf("[MESSAGE] inserting conversation=%d sender=%d\n", conversationID, senderID)

	result, err := tx.Exec(`
		INSERT INTO MESSAGES (
			conversation_id,
			sender_id,
			text
		)
		VALUES (?, ?, ?)
	`,
		conversationID,
		senderID,
		data.Text,
	)
	if err != nil {
		fmt.Println("[MESSAGE] insert error:", err)
		return
	}

	messageID, _ := result.LastInsertId()

	fmt.Printf("[MESSAGE] created id=%d\n", messageID)

	// -------------------------
	// Update conversation preview
	// -------------------------
	fmt.Printf("[CONVERSATION] updating preview id=%d\n", conversationID)

	_, err = tx.Exec(`
		UPDATE CONVERSATIONS
		SET
			last_message = ?,
			last_message_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		data.Text,
		conversationID,
	)
	if err != nil {
		fmt.Println("[CONVERSATION] update preview error:", err)
		return
	}

	// -------------------------
	// Commit
	// -------------------------
	if err := tx.Commit(); err != nil {
		fmt.Println("[DB] commit error:", err)
		return
	}

	fmt.Printf(
		"[SUCCESS] conversation=%d message=%d sender=%d receiver=%d\n",
		conversationID,
		messageID,
		senderID,
		data.ReceiverID,
	)

	fmt.Println("========== send the socket events ==========")
	NotifyUser(
		strconv.Itoa(data.ReceiverID),
		"new_message",
		map[string]interface{}{
			"conversation_id": conversationID,
			"message_id":      messageID,
			"sender_id":       senderID,
			"text":            data.Text,
		},
	)
	fmt.Println("========== CREATE MESSAGE END ==========")

	// Optionally, notify the sender's own connection too (e.g. to sync across their other devices/tabs):
	// ws.NotifyUser(
	// 	strconv.Itoa(senderID),
	// 	"message_sent",
	// 	map[string]interface{}{
	// 		"conversation_id": conversationID,
	// 		"message_id":      messageID,
	// 	},
	// )
}
