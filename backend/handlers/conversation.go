package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"forum/database"
	"forum/utilities"
	"forum/ws"
)

type SendMessageRequest struct {
	ReceiverID     int    `json:"receiver_id"`
	Text           string `json:"text"`
	ConversationID *int   `json:"conversation_id"`
}

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

type Profile struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

type ConversationPreview struct {
	ConversationID *int    `json:"conversationId"`
	Date           *string `json:"date"`
	LastMessage    *string `json:"lastMessage"`
	Status         string  `json:"status"`

	LastSeen    *string `json:"lastSeen"`
	UnreadCount int     `json:"unreadCount"`
	LastSender  string  `json:"lastSender"`
}

type UserFeedItem struct {
	Profile      Profile             `json:"profile"`
	Conversation ConversationPreview `json:"conversation"`
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("========== SEND MESSAGE START ==========")

	// -------------------------
	// Get sender from session
	// -------------------------
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("[AUTH] missing session cookie")
		utilities.WriteJSON(w, 401, "unauthorized", nil)
		return
	}

	senderID, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		fmt.Println("[AUTH] invalid session:", err)
		utilities.WriteJSON(w, 401, "unauthorized", nil)
		return
	}

	fmt.Printf("[AUTH] sender=%d\n", senderID)

	// -------------------------
	// Decode request
	// -------------------------
	var req SendMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("[REQUEST] decode error:", err)
		utilities.WriteJSON(w, 400, "invalid request body", nil)
		return
	}

	fmt.Printf(
		"[REQUEST] receiver=%d text=%q conversation_id=%v\n",
		req.ReceiverID,
		req.Text,
		req.ConversationID,
	)

	// -------------------------
	// Validate
	// -------------------------
	if req.ReceiverID == 0 || req.Text == "" {
		fmt.Println("[VALIDATION] missing fields")
		utilities.WriteJSON(w, 400, "missing fields", nil)
		return
	}

	if senderID == req.ReceiverID {
		fmt.Println("[VALIDATION] user tried to message himself")
		utilities.WriteJSON(w, 400, "cannot send message to yourself", nil)
		return
	}

	// -------------------------
	// Normalize pair
	// -------------------------
	user1 := senderID
	user2 := req.ReceiverID

	if user1 > user2 {
		user1, user2 = user2, user1
	}

	fmt.Printf(
		"[CONVERSATION] normalized pair=(%d,%d)\n",
		user1,
		user2,
	)

	// -------------------------
	// Start transaction
	// -------------------------
	tx, err := database.Database.Begin()
	if err != nil {
		fmt.Println("[DB] begin transaction error:", err)
		utilities.WriteJSON(w, 500, "db error", nil)
		return
	}

	defer tx.Rollback()

	var conversationID int

	// -------------------------
	// CASE 1:
	// conversation_id provided
	// -------------------------
	if req.ConversationID != nil {

		conversationID = *req.ConversationID

		fmt.Printf(
			"[CONVERSATION] validating conversation_id=%d\n",
			conversationID,
		)

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
			fmt.Println(
				"[CONVERSATION] invalid conversation:",
				err,
			)

			utilities.WriteJSON(
				w,
				400,
				"invalid conversation",
				nil,
			)
			return
		}

		fmt.Printf(
			"[CONVERSATION] validated id=%d\n",
			exists,
		)

	} else {

		// -------------------------
		// CASE 2:
		// Find or create conversation
		// -------------------------
		fmt.Printf(
			"[CONVERSATION] searching (%d,%d)\n",
			user1,
			user2,
		)

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

			fmt.Printf(
				"[CONVERSATION] not found, creating (%d,%d)\n",
				user1,
				user2,
			)

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
				fmt.Println(
					"[CONVERSATION] create error:",
					err,
				)

				utilities.WriteJSON(
					w,
					500,
					"failed to create conversation",
					nil,
				)
				return
			}

			id, err := res.LastInsertId()
			if err != nil {
				fmt.Println(
					"[CONVERSATION] last insert id error:",
					err,
				)

				utilities.WriteJSON(
					w,
					500,
					"failed to create conversation",
					nil,
				)
				return
			}

			conversationID = int(id)

			fmt.Printf(
				"[CONVERSATION] created id=%d\n",
				conversationID,
			)

		} else if err != nil {

			fmt.Println(
				"[CONVERSATION] lookup error:",
				err,
			)

			utilities.WriteJSON(
				w,
				500,
				"db error",
				nil,
			)
			return

		} else {
			fmt.Printf(
				"[CONVERSATION] found id=%d\n",
				conversationID,
			)
		}
	}

	// -------------------------
	// Insert message
	// -------------------------
	fmt.Printf(
		"[MESSAGE] inserting conversation=%d sender=%d\n",
		conversationID,
		senderID,
	)

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
		req.Text,
	)
	if err != nil {
		fmt.Println(
			"[MESSAGE] insert error:",
			err,
		)

		utilities.WriteJSON(
			w,
			500,
			"failed to send message",
			nil,
		)
		return
	}

	messageID, _ := result.LastInsertId()

	fmt.Printf(
		"[MESSAGE] created id=%d\n",
		messageID,
	)

	// -------------------------
	// Update conversation preview
	// -------------------------
	fmt.Printf(
		"[CONVERSATION] updating preview id=%d\n",
		conversationID,
	)

	_, err = tx.Exec(`
		UPDATE CONVERSATIONS
		SET
			last_message = ?,
			last_message_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		req.Text,
		conversationID,
	)
	if err != nil {
		fmt.Println(
			"[CONVERSATION] update preview error:",
			err,
		)

		utilities.WriteJSON(
			w,
			500,
			"failed to update conversation",
			nil,
		)
		return
	}

	// -------------------------
	// Commit
	// -------------------------
	if err := tx.Commit(); err != nil {
		fmt.Println("[DB] commit error:", err)

		utilities.WriteJSON(
			w,
			500,
			"commit failed",
			nil,
		)
		return
	}

	fmt.Printf(
		"[SUCCESS] conversation=%d message=%d sender=%d receiver=%d\n",
		conversationID,
		messageID,
		senderID,
		req.ReceiverID,
	)

	fmt.Println("========== send the socket events ==========")
	ws.NotifyUser(
		strconv.Itoa(req.ReceiverID),
		"new_message",
		map[string]interface{}{
			"conversation_id": conversationID,
			"message_id":      messageID,
			"sender_id":       senderID,
			"text":            req.Text,
		},
	)
	fmt.Println("send:", senderID)
	fmt.Println("recive ID:", req.ReceiverID)

	ws.NotifyUser(
		strconv.Itoa(senderID),
		"new_message",
		map[string]interface{}{
			"conversation_id": conversationID,
			"message_id":      messageID,
			"sender_id":       req.ReceiverID,
			"text":            req.Text,
		},
	)

	fmt.Println("========== SEND MESSAGE END ==========")

	utilities.WriteJSON(
		w,
		200,
		"message sent success",
		map[string]interface{}{
			"conversation_id": conversationID,
			"message_id":      messageID,
		},
	)
}

// get all users to show in the UI the first 30 and add throttle to add more 30 by 30   (?offset=10&limit=10)
// rule of sorting 1 for last conversation then alphabitique
func GetConversation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start get users")

	cookie, _ := r.Cookie("session_id")

	userId, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.WriteJSON(w, 405, "not authorized", nil)
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 30
	}
	if limit > 30 {
		limit = 30
	}

	items := []UserFeedItem{}

	// =========================
	// 1. USERS WITH CONVERSATION
	// =========================
	rows, err := database.Database.Query(`
		SELECT 
		    c.id,
			u.id, u.nickname, u.firstname, u.lastname, u.age, u.gender,
			u.last_seen,
			c.last_message,
			c.last_message_at
		FROM USERS u
		JOIN CONVERSATIONS c
			ON (
				(c.user1_id = ? AND c.user2_id = u.id)
				OR
				(c.user2_id = ? AND c.user1_id = u.id)
			)
		WHERE u.id != ?
		ORDER BY c.last_message_at DESC
		LIMIT ? OFFSET ?;
	`, userId, userId, userId, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			convID   int
			u        Profile
			lastMsg  sql.NullString
			lastDate sql.NullString
			lastSeen sql.NullString
		)

		err := rows.Scan(
			&convID,
			&u.ID,
			&u.Nickname,
			&u.Firstname,
			&u.Lastname,
			&u.Age,
			&u.Gender,
			&lastSeen,
			&lastMsg,
			&lastDate,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// =========================
		// DEFAULT LAST SEEN
		// =========================
		lastSeenStr := "24/06/2026"
		if lastSeen.Valid {
			lastSeenStr = lastSeen.String
		}

		// =========================
		// UNREAD MESSAGES
		// =========================
		var unreadCount int
		_ = database.Database.QueryRow(`
			SELECT COUNT(*)
			FROM MESSAGES
			WHERE conversation_id = ?
			AND sender_id != ?
			AND is_read = 0
		`, convID, userId).Scan(&unreadCount)

		// =========================
		// WHO SENT LAST MESSAGE
		// =========================
		lastSender := "them"
		if lastMsg.Valid {
			// optional: you can improve this later with sender_id in messages
			lastSender = "unknown"
		}

		var msgPtr *string
		if lastMsg.Valid {
			msgPtr = &lastMsg.String
		}

		var datePtr *string
		if lastDate.Valid {
			datePtr = &lastDate.String
		}

		items = append(items, UserFeedItem{
			Profile: u,
			Conversation: ConversationPreview{
				ConversationID: &convID,
				Date:           datePtr,
				LastMessage:    msgPtr,
				Status:         "active",

				// NEW FIELDS YOU SHOULD ADD IN STRUCT
				LastSeen:    &lastSeenStr,
				UnreadCount: unreadCount,
				LastSender:  lastSender,
			},
		})
	}

	// =========================
	// 2. USERS WITHOUT CONVERSATION
	// =========================
	rows2, err := database.Database.Query(`
		SELECT u.id, u.nickname, u.firstname, u.lastname, u.age, u.gender
		FROM USERS u
		WHERE u.id != ?
		AND u.id NOT IN (
			SELECT 
				CASE 
					WHEN user1_id = ? THEN user2_id
					ELSE user1_id
				END
			FROM CONVERSATIONS
			WHERE user1_id = ? OR user2_id = ?
		)
		ORDER BY u.nickname COLLATE NOCASE ASC
		LIMIT ?;
	`, userId, userId, userId, userId, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows2.Close()

	for rows2.Next() {
		var u Profile

		err := rows2.Scan(
			&u.ID,
			&u.Nickname,
			&u.Firstname,
			&u.Lastname,
			&u.Age,
			&u.Gender,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		items = append(items, UserFeedItem{
			Profile: u,
			Conversation: ConversationPreview{
				ConversationID: nil,
				Date:           nil,
				LastMessage:    nil,
				Status:         "new",

				LastSeen:    nil,
				UnreadCount: 0,
				LastSender:  "",
			},
		})
	}

	utilities.WriteJSON(w, 200, "ok", items)
}

func GetConversationByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("========== GET CONVERSATION BY ID ==========")

	// -------------------------
	// AUTH
	// -------------------------
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.WriteJSON(w, 401, "unauthorized", nil)
		return
	}

	userID, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.WriteJSON(w, 401, "unauthorized", nil)
		return
	}

	// -------------------------
	// GET conversation ID
	// -------------------------
	idStr := r.PathValue("convID")
	conversationID, err := strconv.Atoi(idStr)
	if err != nil {
		utilities.WriteJSON(w, 400, "invalid conversation id", nil)
		return
	}

	// -------------------------
	// OFFSET + LIMIT
	// -------------------------
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	fmt.Println("conversation:", conversationID, "offset:", offset, "limit:", limit)

	// -------------------------
	// VERIFY USER BELONGS TO CONVERSATION
	// -------------------------
	var convID int
	err = database.Database.QueryRow(`
		SELECT id
		FROM CONVERSATIONS
		WHERE id = ?
		AND (user1_id = ? OR user2_id = ?)
	`, conversationID, userID, userID).Scan(&convID)

	if err == sql.ErrNoRows {
		utilities.WriteJSON(w, 403, "not allowed", nil)
		return
	}
	if err != nil {
		utilities.WriteJSON(w, 500, "db error", nil)
		return
	}

	// -------------------------
	// GET MESSAGES (PAGINATED)
	// -------------------------
	rows, err := database.Database.Query(`
		SELECT id, sender_id, text, created_at
		FROM MESSAGES
		WHERE conversation_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, conversationID, limit, offset)
	if err != nil {
		utilities.WriteJSON(w, 500, "db error", nil)
		return
	}
	defer rows.Close()

	type Message struct {
		ID        int    `json:"id"`
		SenderID  int    `json:"sender_id"`
		Text      string `json:"text"`
		CreatedAt string `json:"created_at"`
	}

	messages := []Message{}

	for rows.Next() {
		var m Message

		err := rows.Scan(
			&m.ID,
			&m.SenderID,
			&m.Text,
			&m.CreatedAt,
		)
		if err != nil {
			utilities.WriteJSON(w, 500, "scan error", nil)
			return
		}

		messages = append(messages, m)
	}

	// -------------------------
	// RESPONSE
	// -------------------------
	utilities.WriteJSON(w, 200, "ok", map[string]interface{}{
		"conversation_id": convID,
		"messages":        messages,
	})
}
