package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/database"
	"forum/models"
	"forum/utilities"
)

// CreateComment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/comments/create" {
		utilities.HandleError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodPost {
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postId := strings.TrimSpace(r.FormValue("postId"))
	text := strings.TrimSpace(r.FormValue("text"))

	if text == "" {
		utilities.HandleError(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	if postId == "" {
		utilities.HandleError(w, http.StatusBadRequest, "Invalid post")
		return
	}
	if len(text) > 1000 {
		utilities.HandleError(w, http.StatusBadRequest, "Comment cannot exceed 1000 characters")
		return
	}

	_, err := strconv.Atoi(postId)
	if err != nil {
		utilities.HandleError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	cookie, _ := r.Cookie("session_id")
	userId, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}
	if _, err = database.Database.Exec(
		"INSERT INTO comments (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)",
		userId,
		postId,
		time.Now(),
		text,
	); err != nil {
		utilities.HandleError(w, http.StatusInternalServerError, "Could not create comment")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CommentResolver
func CommentResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	userId, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.HandleError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := ReactToComment(userId, commentId, 1); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not react to comment")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "dislike":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := ReactToComment(userId, commentId, -1); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not react to comment")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "delete":
		if r.Method != http.MethodDelete {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := DeleteComment(commentId, userId); err != nil {
			utilities.HandleError(w, http.StatusForbidden, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		utilities.HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}
}

// GetCommentsByPost
func GetCommentsByPost(postId int) ([]models.Comment, error) {
	var comments []models.Comment
	rows, err := database.Database.Query(
		"SELECT id, user_id, created_at, text FROM Comments WHERE post_id = ?",
		postId,
	)
	if err != nil {
		return nil, fmt.Errorf("getCommentsByPost error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.Id, &c.UserId, &c.Created_at, &c.Text); err != nil {
			return nil, fmt.Errorf("getCommentsByPost scan error: %v", err)
		}

		// get username
		if err := database.Database.QueryRow(
			"SELECT u.name FROM users u INNER JOIN comments c ON c.user_id = u.id WHERE c.id = ?",
			c.Id,
		).Scan(&c.Username); err != nil {
			return nil, fmt.Errorf("getCommentsByPost username error: %v", err)
		}

		// get timeago
		c.TimeAgo = utilities.TimeAgo(c.Created_at)

		// get reactions
		if c.LikeCount, c.DislikeCount, err = GetReactionsByComment(c.Id); err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCommentsByPost rows error: %v", err)
	}
	return comments, nil
}

// DeleteComment
func DeleteComment(commentId, userId int) error {
	tx, err := database.Database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var dbUserId int
	err = tx.QueryRow("SELECT user_id FROM comments WHERE id = ?", commentId).Scan(&dbUserId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("comment not found")
	}
	if err != nil {
		return err
	}
	if dbUserId != userId {
		return fmt.Errorf("not your comment")
	}

	_, err = tx.Exec("DELETE FROM comments WHERE id = ?", commentId)
	if err != nil {
		return err
	}
	return tx.Commit()
}
