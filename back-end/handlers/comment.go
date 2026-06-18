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
		utilities.WriteJSON(w, http.StatusNotFound, "Page not found", nil)
		return
	}
	fmt.Println("start creating comment")

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}
	type CommentReq struct {
		PostId any    `json:"postId"`
		Text   string `json:"text"`
	}

	// Content type (optional but fine to keep)
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		utilities.WriteJSON(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
		return
	}

	// ✅ REPLACED PART (clean)
	comment, err := utilities.ReadJSONRequest[CommentReq](r)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}
	postId := comment.PostId
	text := comment.Text
	fmt.Println("postId", postId, "texts", text)

	if text == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "Comment cannot be empty", nil)
		return
	}

	if postId == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "Invalid post", nil)
		return
	}
	if len(text) > 1000 {
		utilities.WriteJSON(w, http.StatusBadRequest, "Comment cannot exceed 1000 characters", nil)
		return
	}

	postIntId, err := utilities.ToInt(postId)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "Invalid post ID", nil)
		return
	}

	cookie, _ := r.Cookie("session_id")
	userId, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "Invalid or expired session", nil)
		return
	}
	fmt.Printf("userId=%v (%T)\n", userId, userId)
	fmt.Printf("postId=%v (%T)\n", postIntId, postIntId)
	fmt.Printf("postId=%v (%T)\n", time.Now(), time.Now())
	fmt.Printf("postId=%v (%T)\n", text, text)

	if _, err = database.Database.Exec(
		"INSERT INTO comments (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)",
		userId,
		postIntId,
		time.Now(),
		text,
	); err != nil {
		fmt.Println("errors", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "Could not create comment", nil)
		return
	}

	type Res struct {
		Text      string    `json:"text"`
		PostID    int       `json:"postId"`
		UserID    int       `json:"userId"`
		CreatedAt time.Time `json:"createdAt"`
	}
	res := Res{
		Text:      text,
		PostID:    postIntId,
		UserID:    userId,
		CreatedAt: time.Now(),
	}

	utilities.WriteJSON(w, http.StatusCreated, "message created successfully", res)
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
			"SELECT u.nickname FROM users u INNER JOIN comments c ON c.user_id = u.id WHERE c.id = ?",
			c.Id,
		).Scan(&c.Nickname); err != nil {
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
