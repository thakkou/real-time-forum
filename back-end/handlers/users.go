package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/database"

	"forum/utilities"
)

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
	Date        *string `json:"date,omitempty"`
	LastMessage *string `json:"lastMessage,omitempty"`
	Status      string  `json:"status"` // "new" or "active"
}

type UserFeedItem struct {
	Profile      Profile             `json:"profile"`
	Conversation ConversationPreview `json:"conversation"`
}

// get all users to show in the UI the first 30 and add throttle to add more 30 by 30   (?offset=10&limit=10)
// rule of sorting 1 for last conversation then alphabitique
func GetUsers(w http.ResponseWriter, r *http.Request) {
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
			u.id, u.nickname, u.firstname, u.lastname, u.age, u.gender,
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
			u        Profile
			lastMsg  sql.NullString
			lastDate sql.NullString
		)

		err := rows.Scan(
			&u.ID,
			&u.Nickname,
			&u.Firstname,
			&u.Lastname,
			&u.Age,
			&u.Gender,
			&lastMsg,
			&lastDate,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
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
				Date:        datePtr,
				LastMessage: msgPtr,
				Status:      "active",
			},
		})
	}

	// If full, return early
	if len(items) >= limit {
		utilities.WriteJSON(w, 200, "ok", items)
		return
	}

	remaining := limit - len(items)

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
	`, userId, userId, userId, userId, remaining)
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
				Date:        nil,
				LastMessage: nil,
				Status:      "new",
			},
		})
	}

	// =========================
	// RESPONSE
	// =========================
	utilities.WriteJSON(w, 200, "ok", items)
}

// get UserProfile
func GetUsersById(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	var user User

	query := `
	SELECT
		id,
		nickname,
		firstname,
		lastname,
		age,
		gender
	FROM USERS
	WHERE id = ?
	`

	err := database.Database.QueryRow(query, id).Scan(
		&user.ID,
		&user.Nickname,
		&user.Firstname,
		&user.Lastname,
		&user.Age,
		&user.Gender,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utilities.WriteJSON(w, 200, "user data  get succes", user)
}
