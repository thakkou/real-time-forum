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
		return
	}
	if r.Method != http.MethodPost {
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	text := strings.TrimSpace(r.FormValue("text"))
	categories := r.Form["categories"]

	if title == "" || text == "" {
		HandleError(w, http.StatusBadRequest, "Title and text cannot be empty")
		return
	}

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

	if err = database.Database.QueryRow(
		"SELECT user_id FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&userId); err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	tx, err := database.Database.Begin()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not create post")
		return
	}
	defer tx.Rollback()

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

	postID, err := result.LastInsertId()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not retrieve post ID")
		return
	}

	for _, categoryName := range categories {
		var categoryID int
		if err := tx.QueryRow(
			"SELECT id FROM category WHERE name = ?",
			categoryName,
		).Scan(&categoryID); err != nil {
			HandleError(w, http.StatusBadRequest, "Invalid category: "+categoryName)
			return
		}

		if _, err = tx.Exec(
			"INSERT INTO post_category (post_id, category_id) VALUES (?, ?)",
			postID,
			categoryID,
		); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not associate categories with post")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not save post")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments/create" {
		return
	}
	if r.Method != http.MethodPost {
		return
	}

	postId := strings.TrimSpace(r.FormValue("postId"))
	text := strings.TrimSpace(r.FormValue("text"))

	if text == "" {
		HandleError(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	if postId == "" {
		HandleError(w, http.StatusBadRequest, "Invalid post")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "You must be logged in to comment")
		return
	}

	var userId int
	if err = database.Database.QueryRow(
		"SELECT user_id FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&userId); err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	if _, err = database.Database.Exec(
		"INSERT INTO comments (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)",
		userId,
		postId,
		time.Now(),
		text,
	); err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not create comment")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

///////////////////////////////////////////////////////////

func PostResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	user, err := getUser(cookie.Value)
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.ReactToPost(user.Id, postId, 1); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not react to post")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "dislike":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.ReactToPost(user.Id, postId, -1); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not react to post")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "delete":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.DeletePost(postId, user.Id); err != nil {
			HandleError(w, http.StatusForbidden, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}
}

func CommentResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	user, err := getUser(cookie.Value)
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.ReactToComment(user.Id, commentId, 1); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not react to comment")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "dislike":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.ReactToComment(user.Id, commentId, -1); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not react to comment")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "delete":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.DeleteComment(commentId, user.Id); err != nil {
			HandleError(w, http.StatusForbidden, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}
}