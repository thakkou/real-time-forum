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
	"forum/ws"
)

// =========================
// CORE POST ENRICHMENT
// =========================

func enrichPost(p *models.Post, userId int) error {
	p.TimeAgo = utilities.TimeAgo(p.Created_at)

	// =========================
	// USER INFO
	// =========================
	if err := database.Database.QueryRow(
		"SELECT nickname FROM users WHERE id = ?",
		p.UserId,
	).Scan(&p.Nickname); err != nil {
		return err
	}

	// =========================
	// REACTIONS COUNT
	// =========================
	var err error
	p.LikeCount, p.DislikeCount, err = GetReactionsByPost(p.Id)
	if err != nil {
		return err
	}

	// =========================
	// USER REACTION (IMPORTANT)
	// =========================
	var isLike int

	err = database.Database.QueryRow(`
		SELECT is_like
		FROM POST_REACTIONS
		WHERE user_id = ? AND post_id = ?
	`, userId, p.Id).Scan(&isLike)

	if err == sql.ErrNoRows {
		p.IsLiked = 0 // no reaction
	} else if err != nil {
		return err
	} else {
		p.IsLiked = isLike // 1 or -1
	}

	// =========================
	// CATEGORIES
	// =========================
	p.Categories, err = GetCategoriesByPost(p.Id)
	return err
}

func enrichPostWithComments(p *models.Post, userId int) error {
	if err := enrichPost(p, userId); err != nil {
		return err
	}

	comments, err := GetCommentsByPost(p.Id)
	if err != nil {
		return err
	}

	p.Comments = comments
	return nil
}

func scanPost(row *sql.Rows) (models.Post, error) {
	var p models.Post
	err := row.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text)
	return p, err
}

// =========================
// CREATE POST
// =========================

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/posts/create" {
		utilities.WriteJSON(w, http.StatusNotFound, "Page Not Found", nil)
		return
	}

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	type Post struct {
		Title      string   `json:"title"`
		Text       string   `json:"text"`
		Categories []string `json:"categories"`
	}

	post, err := utilities.ReadJSONRequest[Post](r)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if post.Title == "" || post.Text == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "Title and text cannot be empty", nil)
		return
	}

	if len(post.Title) > 255 || len(post.Text) > 1000 {
		utilities.WriteJSON(w, http.StatusBadRequest, "Title too long", nil)
		return
	}

	if len(post.Categories) == 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "At least one category required", nil)
		return
	}

	cookie, _ := r.Cookie("session_id")
	userId, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid session", nil)
		return
	}

	tx, err := database.Database.Begin()
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "DB error", nil)
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO posts (user_id, created_at, title, text) VALUES (?, ?, ?, ?)",
		userId, time.Now(), post.Title, post.Text,
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Create failed", nil)
		return
	}

	postID, _ := res.LastInsertId()

	for _, c := range post.Categories {
		var catID int
		if err := tx.QueryRow("SELECT id FROM category WHERE name = ?", c).Scan(&catID); err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "Invalid category: "+c, nil)
			return
		}

		_, err := tx.Exec(
			"INSERT INTO post_category (post_id, category_id) VALUES (?, ?)",
			postID, catID,
		)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "category link failed", nil)
			return
		}
	}

	tx.Commit()
	createdPost, err := GetPost(int(postID))

	go ws.BroadcastExcept(strconv.Itoa(userId), "new_post", createdPost)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch created post", nil)
		return
	}

	utilities.WriteJSON(w, 200, "post created successfully", createdPost)
}

// =========================
// POST RESOLVER
// =========================

func PostResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	userId, _ := utilities.GetUserIDFromCookie(cookie.Value)
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid post id", nil)
		return
	}

	switch endpoint {

	// =========================
	// LIKE
	// =========================
	case "like":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		_ = ReactToPost(userId, postId, 1)

		likes, dislikes, _ := GetReactionsByPost(postId)

		utilities.WriteJSON(w, 200, "liked", map[string]any{
			"postId":       postId,
			"likes":        likes,
			"dislikes":     dislikes,
			"userReaction": "like",
		})

	// =========================
	// DISLIKE
	// =========================
	case "dislike":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		_ = ReactToPost(userId, postId, -1)

		likes, dislikes, _ := GetReactionsByPost(postId)

		utilities.WriteJSON(w, 200, "disliked", map[string]any{
			"postId":       postId,
			"likes":        likes,
			"dislikes":     dislikes,
			"userReaction": "dislike",
		})

	// =========================
	// DELETE
	// =========================
	case "delete":
		if r.Method != http.MethodDelete {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		if err := DeletePost(postId, userId); err != nil {
			utilities.WriteJSON(w, 403, err.Error(), nil)
			return
		}

		utilities.WriteJSON(w, 200, "deleted", map[string]any{
			"postId": postId,
		})

	default:
		utilities.WriteJSON(w, 404, "unknown endpoint", nil)
	}
}

// =========================
// GET POSTS
// =========================

func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, 405, "Method not allowed", nil)
		return
	}

	_ = r.ParseForm()

	categories := r.Form["categories"]
	liked := r.FormValue("my-liked-posts") == "true"
	byMe := r.FormValue("my-creat-posts") == "true"

	limit := 30
	offset := 0

	if l := r.FormValue("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	if o := r.FormValue("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	var userID int
	if cookie, err := r.Cookie("session_id"); err == nil {
		userID, _ = utilities.GetUserIDFromCookie(cookie.Value)
	}

	posts, err := GetFilteredPosts(userID, categories, liked, byMe, limit, offset)
	if err != nil {
		utilities.WriteJSON(w, 500, "error", nil)
		return
	}

	utilities.WriteJSON(w, 200, "ok", posts)
}

// =========================
// SINGLE POST
// =========================

func GetPostById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.WriteJSON(w, 400, "Invalid ID", nil)
		return
	}

	post, err := GetPost(id)
	if err != nil {
		utilities.WriteJSON(w, 404, "Not found", nil)
		return
	}

	comments, _ := GetCommentsByPost(id)
	post.Comments = comments

	utilities.WriteJSON(w, 200, "ok", post)
}

// =========================
// FILTERED POSTS
// =========================

func GetFilteredPosts(
	userID int,
	categories []string,
	likedByMe, postedByMe bool,
	limit int,
	offset int,
) ([]models.Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text
		FROM posts p
		LEFT JOIN post_category pc ON p.id = pc.post_id
		LEFT JOIN category c ON pc.category_id = c.id
	`

	var cond []string
	var args []any

	if len(categories) > 0 {
		ph := []string{}
		for _, c := range categories {
			ph = append(ph, "?")
			args = append(args, c)
		}
		cond = append(cond, "c.name IN ("+strings.Join(ph, ",")+")")
	}

	if postedByMe {
		cond = append(cond, "p.user_id = ?")
		args = append(args, userID)
	}

	if likedByMe {
		query += " JOIN post_reactions pr ON p.id = pr.post_id "
		cond = append(cond, "pr.user_id = ? AND pr.is_like = 1")
		args = append(args, userID)
	}

	if len(cond) > 0 {
		query += " WHERE " + strings.Join(cond, " AND ")
	}

	query += " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.Database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}

		if err := enrichPost(&p, userID); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, rows.Err()
}

// =========================
// DELETE POST
// =========================

func DeletePost(postId, userId int) error {
	tx, err := database.Database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var owner int
	err = tx.QueryRow("SELECT user_id FROM posts WHERE id = ?", postId).Scan(&owner)
	if err == sql.ErrNoRows {
		return fmt.Errorf("post not found")
	}
	if owner != userId {
		return fmt.Errorf("not your post")
	}

	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", postId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// =========================
// SINGLE POST FETCH
// =========================

func GetPost(postID int) (models.Post, error) {
	var p models.Post

	err := database.Database.QueryRow(`
		SELECT id, user_id, created_at, title, text
		FROM posts
		WHERE id = ?
	`, postID).Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text)
	if err != nil {
		return p, err
	}

	if err := enrichPostWithComments(&p, p.UserId); err != nil {
		return p, err
	}

	return p, nil
}

func GetUserReaction(userId, postId int) (int, error) {
	var reaction int

	err := database.Database.QueryRow(`
		SELECT is_like
		FROM POST_REACTIONS
		WHERE user_id = ? AND post_id = ?
	`, userId, postId).Scan(&reaction)

	if err == sql.ErrNoRows {
		return 0, nil // no reaction
	}

	return reaction, err
}
