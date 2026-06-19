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

// get all users to show in the UI the first 30 and add throttle to add more 30 by 30   (?offset=10&limit=10)
// rule of sorting 1 for last conversation then alphabitique
func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start get users")
	cookie, _ := r.Cookie("session_id")

	userId, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.WriteJSON(w, 405, "not othorize", nil)
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
	queryUsers := `
SELECT id, nickname, firstname, lastname, age, gender
FROM USERS
WHERE id != ?
ORDER BY nickname COLLATE NOCASE ASC
LIMIT ? OFFSET ?;`

// i wanna get all queryies

//get all converstation users  

//sort them


	// query := `
	// SELECT
	// 	u.id,
	// 	u.nickname,
	// 	u.firstname,
	// 	u.lastname,
	// 	u.age,
	// 	u.gender
	// FROM USERS u
	// LEFT JOIN CONVERSATIONS c
	// 	ON (
	// 		(c.user1_id = ? AND c.user2_id = u.id)
	// 		OR
	// 		(c.user2_id = ? AND c.user1_id = u.id)
	// 	)
	// WHERE u.id != ?
	// ORDER BY
	// 	CASE
	// 		WHEN c.last_message_at IS NULL THEN 1
	// 		ELSE 0
	// 	END,
	// 	c.last_message_at DESC,
	// 	u.nickname COLLATE NOCASE ASC
	// LIMIT ?
	// OFFSET ?
	// `

	rows, err := database.Database.Query(
		queryTes,
		userId,
		userId,
		userId,
		limit,
		offset,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User

		err := rows.Scan(
			&u.ID,
			&u.Nickname,
			&u.Firstname,
			&u.Lastname,
			&u.Age,
			&u.Gender,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users = append(users, u)
	}
	utilities.WriteJSON(w, 200, "users get succes", users)
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
