package api

import (
	"forum/database"
)

// TODO:
// 1. update ui on reaction
// 2. fix like and dislike bugs in comments !!!

func ReactToPost(userId, postId int, isLike bool) {
	var isLiked bool
	err := database.Database.QueryRow(
		"SELECT is_like FROM post_reactions WHERE user_id = ? AND post_id = ?",
		userId,
		postId,
	).Scan(&isLiked) // convert int to bool ?????????????????
	if err == nil {
		// delete
		_, _ = database.Database.Exec(
			"DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?",
			userId,
			postId,
		)
	}
	if !isLiked && isLike || isLiked && !isLike {
		// change reaction
		_, _ = database.Database.Exec(
			"INSERT INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
			userId,
			postId,
			isLike,
		)
	}
}

func ReactToComment(userId, commentId int, isLike bool) {
	var isLiked bool
	err := database.Database.QueryRow(
		"SELECT is_like FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
		userId,
		commentId,
	).Scan(&isLiked) // convert int to bool ??????????????
	if err == nil {
		// delete
		_, _ = database.Database.Exec(
			"DELETE FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
			userId,
			commentId,
		)
	}
	if !isLiked && isLike || isLiked && !isLike {
		// change reaction
		_, _ = database.Database.Exec(
			"INSERT INTO comment_reactions (user_id, comment_id, is_like) VALUES (?, ?, ?)",
			userId,
			commentId,
			isLike,
		)
	}
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
	getNumOfReactions(1, &like_count)    // likes
	getNumOfReactions(0, &dislike_count) // dislikes
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
	getNumOfReactions(1, &like_count)    // likes
	getNumOfReactions(0, &dislike_count) // dislikes
	return like_count, dislike_count, nil
}
