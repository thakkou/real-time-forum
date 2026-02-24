package handlers

import (
	"net/http"
	"strings"
	"time"

	"forum/database"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/posts/create" {
		// handle error with just a status
		// HandleError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodPost {
		// HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	text := r.FormValue("text")

	// handle empty title or text !!!

	// get user id
	var userId int
	cookie, _ := r.Cookie("session_id")
	err := database.Database.QueryRow(
		"SELECT user_id FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&userId)

	// create post
	_, err = database.Database.Exec(
		"INSERT INTO posts (user_id, created_at, title, text) VALUES (?, ?, ?, ?)",
		userId,
		time.Now(),
		title,
		text,
	)
	// create session if you want to redirect to its page
	if err != nil {
		// log.Println(err.Error())
		// HandleError(w, http.StatusInternalServerError, "Could not create account")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
