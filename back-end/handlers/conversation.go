package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"forum/database"
	"forum/utilities"
)

type SendMessageRequest struct {
	ReceiverID     int    `json:"receiver_id"`
	Text           string `json:"text"`
	ConversationID *int   `json:"conversation_id"`
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
