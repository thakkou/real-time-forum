package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/database"
	api "forum/forum-api"
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
	categories := r.Form["categories"] // handle empty categories ?!

	// Validate that at least one category is selected
	if len(categories) == 0 {
		HandleError(w, http.StatusBadRequest, "At least one category must be selected")
		return
	}

	fmt.Println("data", title, text, categories)

	var userId int
	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "You must be logged in to create a post")
		return
	}

	err = database.Database.QueryRow(
		"SELECT user_id FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&userId)
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	// Start a transaction to ensure data consistency
	tx, err := database.Database.Begin()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not create post")
		return
	}
	defer tx.Rollback() // Rollback if anything fails

	// Create post
	result, err := tx.Exec(
		"INSERT INTO posts (user_id, created_at, title, text) VALUES (?, ?, ?, ?)",
		userId,
		time.Now(),
		title,
		text,
	)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not create post")
		return
	}

	// Get the ID of the newly created post
	postID, err := result.LastInsertId()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not retrieve post ID")
		return
	}

	// Validate and insert categories
	for _, categoryName := range categories {
		var categoryID int
		err := tx.QueryRow(
			"SELECT id FROM category WHERE name = ?",
			categoryName,
		).Scan(&categoryID)
		if err != nil {
			// Category doesn't exist
			HandleError(w, http.StatusBadRequest, "Invalid category: "+categoryName)
			return
		}

		// Insert into post_category
		_, err = tx.Exec(
			"INSERT INTO post_category (post_id, category_id) VALUES (?, ?)",
			postID,
			categoryID,
		)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not associate categories with post")
			return
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not save post")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments/create" {
		// handle error with just a status
		// HandleError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodPost {
		// HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postId := r.FormValue("postId")
	text := r.FormValue("text") // trim space ?!

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
		"INSERT INTO comments (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)",
		userId,
		postId,
		time.Now(),
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

///////////////////////////////////////////////////////////

func PostResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")
	cookie, _ := r.Cookie("session_id") // http.ErrNoCookie
	user, _ := getUser(cookie.Value)
	postId, _ := strconv.Atoi(r.PathValue("id"))

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		api.ReactToPost(user.Id, postId, true)

	case "dislike":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		api.ReactToPost(user.Id, postId, false)

		// + case delete
	}
}

func CommentResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")
	cookie, _ := r.Cookie("session_id") // http.ErrNoCookie
	user, _ := getUser(cookie.Value)
	commentId, _ := strconv.Atoi(r.PathValue("id"))

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		api.ReactToComment(user.Id, commentId, true)

	case "dislike":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		api.ReactToComment(user.Id, commentId, false)

		// + case delete
	}
}
