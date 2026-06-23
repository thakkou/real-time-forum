package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"forum/database"

	"forum/utilities"
)

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

func GetUsernameByToken(w http.ResponseWriter, r *http.Request) {
	// errors not checked !!!
	cookie, _ := r.Cookie("session_id")
	var userId, nickname string
	database.Database.QueryRow( // returns error
		"SELECT user_id FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&userId)
	database.Database.QueryRow( // returns error
		"SELECT nickname FROM users WHERE id = ?",
		userId,
	).Scan(&nickname)
	utilities.WriteJSON(w, 200, "success", map[string]any{
		"authenticated": true,
		"nickname":      nickname,
	})
}
