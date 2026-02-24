package database

import (
	"fmt"
	"time"
)

type Post struct {
	Created_at time.Time
	Title      string
	Text       string
}

func CreatePost() {
}

func GetPosts() ([]Post, error) {
	var posts []Post
	rows, err := Database.Query(
		"SELECT created_at, title, text FROM posts",
	)
	defer rows.Close() // release database resources
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Created_at, &p.Title, &p.Text); err != nil {
			return nil, fmt.Errorf("getPosts error: %v", err)
		}
		posts = append(posts, p)
	}
	// Important: Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getPosts error: %v", err)
	}
	return posts, err
}
