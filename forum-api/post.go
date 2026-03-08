package api

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"forum/database"
)

type Post struct {
	Id                      int
	UserId                  int
	Username                string
	Created_at              time.Time
	Title                   string
	Text                    string
	LikeCount, DislikeCount int
	IsLiked                 int // 1:liked, 0:disliked, -1:none
	Comments                []Comment
	Categories              []string // Add this field to store category names
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

		// Get comments for the post
		if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
			return nil, err
		}

		// Get categories for the post
		categories, err := GetCategoriesByPost(p.Id)
		if err != nil {
			return nil, err
		}
		p.Categories = categories

		posts = append(posts, p)
	}

	// Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getPosts error: %v", err)
	}

	return posts, nil
}

// this function is for filtrt posts

func GetFiltrtPOst(userID int, categories []string, likedByMe, postedByMe bool) ([]Post, error) {
	fmt.Println("start filtering")
	fmt.Println("cats", categories)
	fmt.Println("liked by me", likedByMe)
	fmt.Println("posted by me", postedByMe)
	db := database.Database

	query := `
		SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text
		FROM POSTS p
		LEFT JOIN POST_CATEGORY pc ON p.id = pc.post_id
		LEFT JOIN CATEGORY c ON pc.category_id = c.id
	`
	// Conditions slice
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

	// TODO: Implement later
	// // Filter by posts by me
	// if postedByMe {
	// 	conditions = append(conditions, "p.user_id = ?")
	// 	args = append(args, userID)
	// }

	// // Filter by liked by me
	// if likedByMe {
	// 	query += `
	// 		JOIN LIKES l ON p.id = l.post_id
	// 	`
	// 	conditions = append(conditions, "l.user_id = ?")
	// 	args = append(args, userID)
	// }

	// Combine conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY p.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		// Fixed: Use correct field names (Id, UserId, Created_at)
		if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text); err != nil {
			return nil, err
		}

		// Get categories for the post
		categories, err := GetCategoriesByPost(p.Id)
		if err != nil {
			return nil, err
		}
		p.Categories = categories

		// Get comments for the post
		if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
			return nil, err
		}

		posts = append(posts, p)
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
			return nil, fmt.Errorf("GetCategoriesByPost error: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCategoriesByPost error: %v", err)
	}

	return categories, nil
}

// Alternative optimized version that gets all posts with their categories in a single query
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
			return nil, fmt.Errorf("GetPostsOptimized error: %v", err)
		}

		// Parse comma-separated categories
		if categoriesStr != "" {
			p.Categories = strings.Split(categoriesStr, ",")
		} else {
			p.Categories = []string{}
		}

		// Get comments for the post
		if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPostsOptimized error: %v", err)
	}

	return posts, nil
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
