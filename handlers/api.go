package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"forum/database"
	api "forum/forum-api"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/posts/create" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed, "Method not Allowed")
		return
	}

	// check the size of dat entry
	err := r.ParseMultipartForm(20 << 20) // 20 MB
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Image max size is 20Mb")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	text := strings.TrimSpace(r.FormValue("text"))
	categories := r.Form["categories"]

	if title == "" || text == "" {
		HandleError(w, http.StatusBadRequest, "Title and text cannot be empty")
		return
	}
	if len(title) > 255 || len(text) > 1000 {
		HandleError(w, http.StatusBadRequest, "Title cannot exceed 255 characters")
		return
	}

	if len(categories) == 0 {
		HandleError(w, http.StatusBadRequest, "At least one category must be selected")
		return
	}

	// get userId
	cookie, _ := r.Cookie("session_id")
	userId, err := api.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	// add image
	var imageUri string // default empty
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		fmt.Println("No image uploaded, continuing without it")
	} else {
		defer file.Close()
		imageUri, err = SaveImage(file, fileHeader)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not save image")
			return
		}
	}

	tx, err := database.Database.Begin()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not create post")
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO posts (user_id, created_at, title, text, image) VALUES (?, ?, ?, ?, ?)",
		userId,
		time.Now(),
		title,
		text,
		imageUri,
	)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not create post")
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not retrieve post ID")
		return
	}

	for _, categoryName := range categories {
		var categoryID int
		if err := tx.QueryRow(
			"SELECT id FROM category WHERE name = ?",
			categoryName,
		).Scan(&categoryID); err != nil {
			HandleError(w, http.StatusBadRequest, "Invalid category: "+categoryName)
			return
		}

		if _, err = tx.Exec(
			"INSERT INTO post_category (post_id, category_id) VALUES (?, ?)",
			postID,
			categoryID,
		); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not associate categories with post")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not save post")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/comments/create" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postId := strings.TrimSpace(r.FormValue("postId"))
	text := strings.TrimSpace(r.FormValue("text"))

	if text == "" {
		HandleError(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	if postId == "" {
		HandleError(w, http.StatusBadRequest, "Invalid post")
		return
	}
	if len(text) > 1000 {
		HandleError(w, http.StatusBadRequest, "Comment cannot exceed 1000 characters")
		return
	}

	_, err := strconv.Atoi(postId)
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	cookie, _ := r.Cookie("session_id")
	userId, err := api.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}
	if _, err = database.Database.Exec(
		"INSERT INTO comments (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)",
		userId,
		postId,
		time.Now(),
		text,
	); err != nil {
		HandleError(w, http.StatusInternalServerError, "Could not create comment")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

///////////////////////////////////////////////////////////

func PostResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	userId, err := api.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.ReactToPost(userId, postId, 1); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not react to post")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "dislike":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.ReactToPost(userId, postId, -1); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not react to post")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "delete":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.DeletePost(postId, userId); err != nil {
			HandleError(w, http.StatusForbidden, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}
}

func CommentResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	userId, err := api.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		HandleError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.ReactToComment(userId, commentId, 1); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not react to comment")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "dislike":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.ReactToComment(userId, commentId, -1); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not react to comment")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "delete":
		if r.Method != http.MethodPost {
			HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := api.DeleteComment(commentId, userId); err != nil {
			HandleError(w, http.StatusForbidden, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}
}

func SaveImage(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	// Ensure uploads directory exists
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("unable to create uploads directory: %w", err)
	}

	// rand.Seed(time.Now().UnixNano())
	ext := filepath.Ext(fileHeader.Filename) // keep original extension
	newName := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), rand.Intn(10000), ext)
	filePath := filepath.Join("./uploads", newName)

	// Create destination file

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("unable to save file: %w", err)
	}

	// Return the relative URL/path for DB insertion
	return "/uploads/" + newName, nil
}
