package api

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"forum/database"
)

type Post struct {
	Id int
	UserId int
	Username string
	Created_at time.Time
	TimeAgo string
	Title string
	Text string
	LikeCount, DislikeCount int
	IsLiked int // 1:liked, 0:none, -1:disliked
	Comments []Comment
	Categories []string
}

func timeAgo(t time.Time) string {
	d := time.Since(t)

	if d < time.Minute {
		return fmt.Sprintf("%d seconds ago", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	}
	if d < 30*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
	return fmt.Sprintf("%d months ago", int(d.Hours()/(24*30)))
}

func GetPosts() ([]Post, error) {
	var posts []Post
		// here i will check the categories and filters

	// Modified query to include user_id since we need it for categories

	rows, err := database.Database.Query(
		"SELECT id, user_id, created_at, title, text FROM posts ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("getPosts error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text); err != nil {
			return nil, fmt.Errorf("getPosts scan error: %v", err)
		}

		// get username
		if err := database.Database.QueryRow(
			"SELECT u.name FROM users u INNER JOIN posts p ON p.user_id = u.id WHERE p.id = ?",
			p.Id,
		).Scan(&p.Username); err != nil {
			return nil, fmt.Errorf("getPosts username error: %v", err)
		}

		// get timeago
		p.TimeAgo = timeAgo(p.Created_at)

		// get reactions
		if p.LikeCount, p.DislikeCount, err = GetReactionsByPost(p.Id); err != nil {
			return nil, fmt.Errorf("getPosts reactions error: %v", err)
		}

		// Get comments for the post
		if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
			return nil, fmt.Errorf("getPosts comments error: %v", err)
		}

		// Get categories for the post
		if p.Categories, err = GetCategoriesByPost(p.Id); err != nil {
			return nil, fmt.Errorf("getPosts categories error: %v", err)
		}
		posts = append(posts, p)
	}
	// Check for any errors that occurred during iteration

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getPosts rows error: %v", err)
	}

	return posts, nil
}
// this function is for filtrt posts

func GetFiltrtPOst(userID int, categories []string, likedByMe, postedByMe bool) ([]Post, error) {
	db := database.Database

	query := `
		SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text
		FROM POSTS p
		LEFT JOIN POST_CATEGORY pc ON p.id = pc.post_id
		LEFT JOIN CATEGORY c ON pc.category_id = c.id
	`

	conditions := []string{}
	args := []interface{}{}
	// Filter by categories

	if len(categories) > 0 {
		placeholders := []string{}
		for _, cat := range categories {
			placeholders = append(placeholders, "?")
			args = append(args, cat)
		}
		
		conditions = append(conditions, "c.name IN ("+strings.Join(placeholders, ",")+")")
	}
	// Filter posts created by user

	if postedByMe && userID != 0 {
		conditions = append(conditions, "p.user_id = ?")
		args = append(args, userID)
	}

	if likedByMe && userID != 0 {
		query += `
			JOIN POST_REACTIONS pr ON p.id = pr.post_id
		`
		conditions = append(conditions, "pr.user_id = ? AND pr.is_like = 1")
		args = append(args, userID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY p.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetFiltrtPOst query error: %v", err)
	}
	defer rows.Close()

	posts := []Post{}

	for rows.Next() {
		var p Post

		if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text); err != nil {
			return nil, fmt.Errorf("GetFiltrtPOst scan error: %v", err)
		}

		// get username
		if err := db.QueryRow(
			"SELECT name FROM users WHERE id = ?",
			p.UserId,
		).Scan(&p.Username); err != nil {
			return nil, fmt.Errorf("GetFiltrtPOst username error: %v", err)
		}

		// get timeago
		p.TimeAgo = timeAgo(p.Created_at)

		// get reactions
		if p.LikeCount, p.DislikeCount, err = GetReactionsByPost(p.Id); err != nil {
			return nil, fmt.Errorf("GetFiltrtPOst reactions error: %v", err)
		}

		// get comments
		if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
			return nil, fmt.Errorf("GetFiltrtPOst comments error: %v", err)
		}

		// get categories
		if p.Categories, err = GetCategoriesByPost(p.Id); err != nil {
			return nil, fmt.Errorf("GetFiltrtPOst categories error: %v", err)
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFiltrtPOst rows error: %v", err)
	}

	return posts, nil
}
// Helper function to get categories for a specific post

func GetCategoriesByPost(postId int) ([]string, error) {
	var categories []string

	rows, err := database.Database.Query(`
		SELECT c.name
		FROM category c
		JOIN post_category pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
		ORDER BY c.name
	`, postId)
	if err != nil {
		return nil, fmt.Errorf("GetCategoriesByPost error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("GetCategoriesByPost scan error: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCategoriesByPost rows error: %v", err)
	}

	return categories, nil
}

func GetPostsOptimized() ([]Post, error) {
	var posts []Post

	// This query gets all posts with their categories aggregated as a JSON array or comma-separated string
	// Since SQLite doesn't have native JSON functions in older versions, we'll use GROUP_CONCAT
	rows, err := database.Database.Query(`
		SELECT
			p.id,
			p.user_id,
			p.created_at,
			p.title,
			p.text,
			COALESCE(GROUP_CONCAT(c.name, ','), '') as categories
		FROM posts p
		LEFT JOIN post_category pc ON p.id = pc.post_id
		LEFT JOIN category c ON pc.category_id = c.id
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("GetPostsOptimized error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		var categoriesStr string

		if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text, &categoriesStr); err != nil {
			return nil, fmt.Errorf("GetPostsOptimized scan error: %v", err)
		}

		// get username
		if err := database.Database.QueryRow(
			"SELECT name FROM users WHERE id = ?", p.UserId,
		).Scan(&p.Username); err != nil {
			return nil, fmt.Errorf("GetPostsOptimized username error: %v", err)
		}

		// get timeago
		p.TimeAgo = timeAgo(p.Created_at)

		// get reactions
		if p.LikeCount, p.DislikeCount, err = GetReactionsByPost(p.Id); err != nil {
			return nil, fmt.Errorf("GetPostsOptimized reactions error: %v", err)
		}

		// parse categories
		if categoriesStr != "" {
			p.Categories = strings.Split(categoriesStr, ",")
		} else {
			p.Categories = []string{}
		}

		// get comments
		if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
			return nil, fmt.Errorf("GetPostsOptimized comments error: %v", err)
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPostsOptimized rows error: %v", err)
	}

	return posts, nil
}

func CheckLikedPosts(posts []Post, userId int) {
	for i, post := range posts {
		_ = database.Database.QueryRow(
			"SELECT is_like FROM post_reactions WHERE user_id = ? AND post_id = ?",
			userId,
			post.Id,
		).Scan(&posts[i].IsLiked)

		for j, comment := range posts[i].Comments {
			_ = database.Database.QueryRow(
				"SELECT is_like FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
				userId,
				comment.Id,
			).Scan(&posts[i].Comments[j].IsLiked)
		}
	}
}

func DeletePost(postId, userId int) error {
	tx, err := database.Database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var dbUserId int
	err = tx.QueryRow("SELECT user_id FROM posts WHERE id = ?", postId).Scan(&dbUserId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("post not found")
	}
	if err != nil {
		return err
	}
	if dbUserId != userId {
		return fmt.Errorf("not your post")
	}

	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", postId)
	if err != nil {
		return err
	}
	return tx.Commit()
}