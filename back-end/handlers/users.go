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
