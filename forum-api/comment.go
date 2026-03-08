package api

import (
	"database/sql"
	"fmt"
	"time"

	"forum/database"
)

type Comment struct {
	Id                      int
	UserId                  int
	Username                string
	Created_at              time.Time
	Text                    string
	LikeCount, DislikeCount int
	IsLiked                 int // 1:liked, 0:disliked, -1:none
}

func GetCommentsByPost(postId int) ([]Comment, error) {
	var comments []Comment
	rows, err := database.Database.Query(
		"SELECT id, user_id, created_at, text FROM Comments WHERE post_id = ?",
		postId,
	)
	defer rows.Close() // release database resources
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.Id, &c.UserId, &c.Created_at, &c.Text); err != nil {
			return nil, fmt.Errorf("getCommentsByPost error: %v", err)
		}

		// get username
		database.Database.QueryRow(
			"SELECT u.name FROM users u INNER JOIN comments c ON c.user_id = u.id WHERE c.id = ?",
			c.Id,
		).Scan(&c.Username)

		// get reactions
		if c.LikeCount, c.DislikeCount, err = GetReactionsByComment(c.Id); err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}
	// Important: Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCommentsByPost error: %v", err)
	}
	return comments, err
}

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
