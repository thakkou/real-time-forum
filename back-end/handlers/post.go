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

// CreatePost
func CreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("posts")
	if r.URL.Path != "/api/posts/create" {
		utilities.HandleError(w, http.StatusNotFound, "Page Not Found")
		return
	}

	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	// this for code is for image upload i will update it later

	// r.Body = http.MaxBytesReader(w, r.Body, 21<<20) // 21 MB hardcoded

	// 1. Get post creation form data

	type Post struct {
		Title      string   `json:"title"`
		Text       string   `json:"text"`
		Categories []string `json:"categories"`
	}

	// Content type (optional but fine to keep)
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		utilities.WriteJSON(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
		return
	}

	// ✅ REPLACED PART (clean)
	post, err := utilities.ReadJSONRequest[Post](r)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	title := post.Title
	text := post.Text
	categories := post.Categories

	// check the size of data entry
	// const maxImageSize int64 = 20 << 20 // 20 MB

	// err := r.ParseMultipartForm(maxImageSize)
	// ParseMultipartForm sets the in-memory buffer limit.
	// If the file exceeds that limit, Go silently spills the overflow to a temp file on disk.
	// if err != nil {
	// 	utilities.HandleError(w, http.StatusBadRequest, "Image max size is 20Mb.") // 400
	// 	return
	// }

	// 2. Sanitize form data
	if title == "" || text == "" {
		utilities.HandleError(w, http.StatusBadRequest, "Title and text cannot be empty")
		return
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if len(title) > 255 || len(text) > 1000 {
		utilities.HandleError(w, http.StatusBadRequest, "Title cannot exceed 255 characters")
		return
	}
	if len(categories) == 0 {
		utilities.HandleError(w, http.StatusBadRequest, "At least one category must be selected")
		return
	}

	// add image
	// var imageUri string // default empty
	// file, header, err := r.FormFile("image")
	// if err != nil {
	// 	// post without image is handled!
	// 	fmt.Println("No image uploaded, continuing without it")
	// } else {
	// 	defer file.Close()

	// 	// Layer 2: header.Size — fast pre-check, avoids reading the file at all (depends on content-length)
	// 	// not 100% trustworthy (client-declared) but useful to reject obviously large files early
	// 	if header.Size > maxImageSize {
	// 		utilities.HandleError(w, http.StatusBadRequest, "Image max size is 20MB")
	// 		return
	// 	}

	// 	// Layer 3: io.LimitReader — trustworthy precise enforcement on actual bytes
	// 	// reads up to maxFileSize+1 to detect if file exceeds the limit
	// 	limitedReader := io.LimitReader(file, maxImageSize+1)
	// 	fileBytes, err := io.ReadAll(limitedReader)
	// 	if err != nil {
	// 		utilities.HandleError(w, http.StatusInternalServerError, "Internal server error")
	// 		return
	// 	}
	// 	if int64(len(fileBytes)) > maxImageSize {
	// 		utilities.HandleError(w, http.StatusRequestEntityTooLarge, "Image max size is 20MB")
	// 		return
	// 	}

	// 	// Check if file is valid image
	// 	buffer := make([]byte, 512)
	// 	file.Seek(0, 0) // without it, Read may give EOF error
	// 	_, err = file.Read(buffer)
	// 	if err != nil && err != io.EOF {
	// 		utilities.HandleError(w, http.StatusInternalServerError, "Could not save image")
	// 		return
	// 	}
	// 	// Reset file pointer so it can be read again later
	// 	if _, err := file.Seek(0, 0); err != nil {
	// 		utilities.HandleError(w, http.StatusInternalServerError, "Could not save image")
	// 		return
	// 	}
	// 	contentType := http.DetectContentType(buffer)
	// 	if !strings.HasPrefix(contentType, "image/") { // svg not handled: complicated + unsafe xml
	// 		utilities.HandleError(w, http.StatusBadRequest, "Invalid image type")
	// 		return
	// 	}

	// 	imageUri, err = utilities.SaveImage(file, header)
	// 	if err != nil {
	// 		utilities.HandleError(w, http.StatusInternalServerError, "Could not save image")
	// 		return
	// 	}
	// }

	// 3. Get userId
	cookie, _ := r.Cookie("session_id")
	userId, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	tx, err := database.Database.Begin()
	if err != nil {
		utilities.HandleError(w, http.StatusInternalServerError, "Could not create post")
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO posts (user_id, created_at, title, text) VALUES (?, ?, ?, ?)",
		userId,
		time.Now(),
		title,
		text,
		// imageUri,
	)
	if err != nil {
		utilities.HandleError(w, http.StatusInternalServerError, "Could not create post")
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		utilities.HandleError(w, http.StatusInternalServerError, "Could not retrieve post ID")
		return
	}

	for _, categoryName := range categories {
		var categoryID int
		if err := tx.QueryRow(
			"SELECT id FROM category WHERE name = ?",
			categoryName,
		).Scan(&categoryID); err != nil {
			utilities.HandleError(w, http.StatusBadRequest, "Invalid category: "+categoryName)
			return
		}

		if _, err = tx.Exec(
			"INSERT INTO post_category (post_id, category_id) VALUES (?, ?)",
			postID,
			categoryID,
		); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not associate categories with post")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		utilities.HandleError(w, http.StatusInternalServerError, "Could not save post")
		return
	}

	utilities.WriteJSON(w, 200, "post created successfully", map[string]any{
		"post_id":    postID,
		"title":      title,
		"text":       text,
		"categories": categories,
		"user_id":    userId,
	})
}

func PostResolver(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post resolver")
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

	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.HandleError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := ReactToPost(userId, postId, 1); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not react to post")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "dislike":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := ReactToPost(userId, postId, -1); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not react to post")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "delete":
		if r.Method != http.MethodDelete {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		if err := DeletePost(postId, userId); err != nil {
			utilities.HandleError(w, http.StatusForbidden, err.Error())
			return
		}
		utilities.WriteJSON(w, 200, "post deleted successfully", map[string]any{
			"post_id": postId,
			"deleted": true,
		})
	default:
		utilities.HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/posts/getPosts" {
		utilities.WriteJSON(w, http.StatusNotFound, "url not found", nil)
		return
	}
	fmt.Println("posts get ")

	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "Bad request", nil)
		return
	}

	categories := r.Form["categories"]
	isLiked := r.FormValue("my-liked-posts") == "true"
	isByMe := r.FormValue("my-creat-posts") == "true"
	fmt.Println("isliked", isLiked)
	fmt.Println("isByMe", isByMe)
	fmt.Println("categories", categories, len(categories))

	var userID int

	cookie, err := r.Cookie("session_id")
	if err == nil {
		userID, _ = utilities.GetUserIDFromCookie(cookie.Value)
	}

	posts, err := GetFilteredPosts(userID, categories, isLiked, isByMe)
	if err != nil {
		fmt.Println("errors", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "failed to get posts", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "posts fetched successfully", posts)
}

// Returns filtered posts
func GetFilteredPosts(userID int, categories []string, likedByMe, postedByMe bool) ([]models.Post, error) {
	fmt.Println("start filtriing posts", userID, categories, likedByMe, postedByMe)

	if len(categories) == 0 && userID == 0 {
		postes, err := GetAllPosts()
		return postes, err
	}

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

	posts := []models.Post{}

	for rows.Next() {
		var p models.Post

		if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text); err != nil {
			return nil, fmt.Errorf("GetFiltrtPOst scan error: %v", err)
		}

		// get nickname
		if err := db.QueryRow(
			"SELECT nickname FROM users WHERE id = ?",
			p.UserId,
		).Scan(&p.Nickname); err != nil {
			return nil, fmt.Errorf("GetFiltrtPOst nickname  error: %v", err)
		}

		// get timeago
		p.TimeAgo = utilities.TimeAgo(p.Created_at)

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

// GetPostsOptimized
func GetPostsOptimized() ([]models.Post, error) {
	var posts []models.Post

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
		var p models.Post
		var categoriesStr string

		if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text, &categoriesStr); err != nil {
			return nil, fmt.Errorf("GetPostsOptimized scan error: %v", err)
		}

		// get username
		if err := database.Database.QueryRow(
			"SELECT name FROM users WHERE id = ?", p.UserId,
		).Scan(&p.Nickname); err != nil {
			return nil, fmt.Errorf("GetPostsOptimized username error: %v", err)
		}

		// get timeago
		p.TimeAgo = utilities.TimeAgo(p.Created_at)

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

// CheckLikedPosts
func CheckLikedPosts(posts []models.Post, userId int) {
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

// DeletePost
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

func GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	// here i will check the categories and filters

	// Modified query to include user_id since we need it for categories

	rows, err := database.Database.Query(
		"SELECT id, user_id, created_at, title, text  FROM posts ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("getPosts error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text); err != nil {
			return nil, fmt.Errorf("getPosts scan error: %v", err)
		}

		// get username
		if err := database.Database.QueryRow(
			"SELECT u.name FROM users u INNER JOIN posts p ON p.user_id = u.id WHERE p.id = ?",
			p.Id,
		).Scan(&p.Nickname); err != nil {
			return nil, fmt.Errorf("getPosts username error: %v", err)
		}

		// get timeago
		p.TimeAgo = utilities.TimeAgo(p.Created_at)

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
