package api

import (
	"fmt"

	"forum/database"
)

func ReactToPost(userId, postId int, isLikeInt int) error {
	var isLikedInt int
	err := database.Database.QueryRow(
		"SELECT is_like FROM post_reactions WHERE user_id = ? AND post_id = ?",
		userId,
		postId,
	).Scan(&isLikedInt)

	if err == nil {
		// delete previous reaction
		if _, err := database.Database.Exec(
			"DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?",
			userId,
			postId,
		); err != nil {
			return fmt.Errorf("ReactToPost delete error: %v", err)
		}
	}

	isLike := isLikeInt == 1
	isLiked := isLikedInt == 1
	if isLike && (err != nil || !isLiked) ||
		!isLike && (err != nil || isLiked) {
		if _, err := database.Database.Exec(
			"INSERT INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
			userId,
			postId,
			isLikeInt,
		); err != nil {
			return fmt.Errorf("ReactToPost insert error: %v", err)
		}
	}

	return nil
}

func ReactToComment(userId, commentId int, isLikeInt int) error {
	var isLiked bool
	err := database.Database.QueryRow(
		"SELECT is_like FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
		userId,
		commentId,
	).Scan(&isLiked)

	if err == nil {
		if _, err := database.Database.Exec(
			"DELETE FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
			userId,
			commentId,
		); err != nil {
			return fmt.Errorf("ReactToComment delete error: %v", err)
		}
	}

	isLike := isLikeInt == 1
	if isLike && (err != nil || !isLiked) ||
		!isLike && (err != nil || isLiked) {
		if _, err := database.Database.Exec(
			"INSERT INTO comment_reactions (user_id, comment_id, is_like) VALUES (?, ?, ?)",
			userId,
			commentId,
			isLikeInt,
		); err != nil {
			return fmt.Errorf("ReactToComment insert error: %v", err)
		}
	}

	return nil
}

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
		return 0, 0, fmt.Errorf("GetReactionsByPost likes error: %v", err)
	}
	if err := getNumOfReactions(-1, &dislike_count); err != nil {
		return 0, 0, fmt.Errorf("GetReactionsByPost dislikes error: %v", err)
	}
	return like_count, dislike_count, nil
}

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
		return 0, 0, fmt.Errorf("GetReactionsByComment likes error: %v", err)
	}
	if err := getNumOfReactions(-1, &dislike_count); err != nil {
		return 0, 0, fmt.Errorf("GetReactionsByComment dislikes error: %v", err)
	}
	return like_count, dislike_count, nil
}