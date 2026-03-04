package api

import (
	"database/sql"
	"fmt"
	"time"

	"forum/database"
)

type Post struct {
	Id                      int
	Username                string
	Created_at              time.Time
	Title                   string
	Text                    string
	LikeCount, DislikeCount int
	IsLiked                 int // 1:liked, 0:disliked, -1:none
	Comments                []Comment
}

func CreatePost() {
}

func GetPosts() ([]Post, error) {
	var posts []Post
	rows, err := database.Database.Query(
		"SELECT id, created_at, title, text FROM posts",
	)
	defer rows.Close() // release database resources
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.Created_at, &p.Title, &p.Text); err != nil {
			return nil, fmt.Errorf("getPosts error: %v", err)
		}

		// get username
		err := database.Database.QueryRow(
			"SELECT u.name FROM users u INNER JOIN posts p ON p.user_id = u.id WHERE p.id = ?",
			p.Id,
		).Scan(&p.Username)

		// get reactions
		if p.LikeCount, p.DislikeCount, err = GetReactionsByPost(p.Id); err != nil {
			return nil, err
		}

		// get comments
		if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	// Important: Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getPosts error: %v", err)
	}
	return posts, err
}

func CheckLikedPosts(posts []Post, userId int) {
	// check if posts are liked
	for i, post := range posts {
		err := database.Database.QueryRow(
			"SELECT is_like FROM post_reactions WHERE user_id = ? AND post_id = ?",
			userId,
			post.Id,
		).Scan(&posts[i].IsLiked)
		if err == sql.ErrNoRows {
			posts[i].IsLiked = -1
		}

		// check if comments are liked
		for j, comment := range posts[i].Comments {
			err := database.Database.QueryRow(
				"SELECT is_like FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
				userId,
				comment.Id,
			).Scan(&posts[i].Comments[j].IsLiked)
			if err == sql.ErrNoRows {
				posts[i].Comments[j].IsLiked = -1
			}
		}
	}
}
