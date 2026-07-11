package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/database"
)

func ReactToPost(userId, postId int, isLikeInt int) (int, error) {
	// Validate reaction
	if isLikeInt != 1 && isLikeInt != -1 {
		return http.StatusBadRequest, fmt.Errorf("invalid reaction")
	}

	// Check if post exists
	var exists int
	err := database.Database.QueryRow(
		"SELECT id FROM posts WHERE id = ?",
		postId,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("post not found")
	}

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("checking post existence: %w", err)
	}

	// Check existing reaction
	var oldReaction int
	err = database.Database.QueryRow(
		"SELECT is_like FROM post_reactions WHERE user_id = ? AND post_id = ?",
		userId,
		postId,
	).Scan(&oldReaction)

	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, fmt.Errorf("checking reaction: %w", err)
	}

	// Existing reaction
	if err == nil {

		// Remove reaction if same
		if oldReaction == isLikeInt {
			_, err = database.Database.Exec(
				"DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?",
				userId,
				postId,
			)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			return http.StatusOK, nil
		}

		// Update reaction
		_, err = database.Database.Exec(
			"UPDATE post_reactions SET is_like = ? WHERE user_id = ? AND post_id = ?",
			isLikeInt,
			userId,
			postId,
		)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, nil
	}

	// Create reaction
	_, err = database.Database.Exec(
		"INSERT INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
		userId,
		postId,
		isLikeInt,
	)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

// ReactToComment
func ReactToComment(userId, commentId int, isLikeInt int) (int, error) {
	// Validate reaction
	if isLikeInt != 1 && isLikeInt != -1 {
		return http.StatusBadRequest, fmt.Errorf("invalid reaction")
	}

	// Check if comment exists
	var exists int
	err := database.Database.QueryRow(
		"SELECT id FROM comments WHERE id = ?",
		commentId,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("comment not found")
	}

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("checking comment existence: %w", err)
	}

	// Check existing reaction
	var oldReaction int
	err = database.Database.QueryRow(
		"SELECT is_like FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
		userId,
		commentId,
	).Scan(&oldReaction)

	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, fmt.Errorf("checking reaction: %w", err)
	}

	// Existing reaction
	if err == nil {

		// Same reaction -> remove it
		if oldReaction == isLikeInt {
			_, err = database.Database.Exec(
				"DELETE FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
				userId,
				commentId,
			)
			if err != nil {
				return http.StatusInternalServerError, fmt.Errorf("delete reaction: %w", err)
			}

			return http.StatusOK, nil
		}

		// Different reaction -> update it
		_, err = database.Database.Exec(
			"UPDATE comment_reactions SET is_like = ? WHERE user_id = ? AND comment_id = ?",
			isLikeInt,
			userId,
			commentId,
		)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("update reaction: %w", err)
		}

		return http.StatusOK, nil
	}

	// No reaction -> insert
	_, err = database.Database.Exec(
		"INSERT INTO comment_reactions (user_id, comment_id, is_like) VALUES (?, ?, ?)",
		userId,
		commentId,
		isLikeInt,
	)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("insert reaction: %w", err)
	}

	return http.StatusCreated, nil
}

// GetReactionsByPost
func GetReactionsByPost(postId int) (int, int, error) {
	var like_count, dislike_count int
	getNumOfReactions := func(is_like int, n *int) error {
		return database.Database.QueryRow(
			"SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND is_like = ?",
			postId,
			is_like,
		).Scan(n)
	}
	if err := getNumOfReactions(1, &like_count); err != nil {
		return 0, 0, fmt.Errorf("GetReactionsByPost likes error: %w", err)
	}
	if err := getNumOfReactions(-1, &dislike_count); err != nil {
		return 0, 0, fmt.Errorf("GetReactionsByPost dislikes error: %w", err)
	}
	return like_count, dislike_count, nil
}

func GetUserCommentReaction(userId, commentId int) (string, error) {
	var isLike int

	err := database.Database.QueryRow(`
		SELECT is_like
		FROM comment_reactions
		WHERE user_id = ? AND comment_id = ?
	`, userId, commentId).Scan(&isLike)

	if err == sql.ErrNoRows {
		return "none", nil
	}

	if err != nil {
		return "", err
	}

	if isLike == 1 {
		return "like", nil
	}

	return "dislike", nil
}

// GetReactionsByComment
func GetReactionsByComment(commentId int) (int, int, error) {
	var like_count, dislike_count int
	getNumOfReactions := func(is_like int, n *int) error {
		return database.Database.QueryRow(
			"SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND is_like = ?",
			commentId,
			is_like,
		).Scan(n)
	}
	if err := getNumOfReactions(1, &like_count); err != nil {
		return 0, 0, fmt.Errorf("GetReactionsByComment likes error: %w", err)
	}
	if err := getNumOfReactions(-1, &dislike_count); err != nil {
		return 0, 0, fmt.Errorf("GetReactionsByComment dislikes error: %w", err)
	}
	return like_count, dislike_count, nil
}
