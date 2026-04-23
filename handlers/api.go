package handlers

import (
	"fmt"
	"io"
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
	if r.URL.Path != "/api/posts/create" {
		utilities.HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodPost {
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// check the size of data entry
	const maxImageSize int64 = 20 << 20 // 20 MB

	// Layer 1: MaxBytesReader — cuts the connection early at network level
	// slightly padded to account for multipart boundaries and form fields overhead
	r.Body = http.MaxBytesReader(w, r.Body, 21<<20) // 21 MB

	err := r.ParseMultipartForm(maxImageSize)
	// ParseMultipartForm sets the in-memory buffer limit.
	// If the file exceeds that limit, Go silently spills the overflow to a temp file on disk.
	if err != nil {
		utilities.HandleError(w, http.StatusBadRequest, "Image max size is 20Mb")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	text := strings.TrimSpace(r.FormValue("text"))
	categories := r.Form["categories"]

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

	// get userId
	cookie, _ := r.Cookie("session_id")
	userId, err := models.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	// add image
	var imageUri string // default empty
	file, header, err := r.FormFile("image")
	if err != nil {
		// post without image is handled!
		fmt.Println("No image uploaded, continuing without it")
	} else {
		defer file.Close()

		// Layer 2: header.Size — fast pre-check, avoids reading the file at all (depends on content-length)
		// not 100% trustworthy (client-declared) but useful to reject obviously large files early
		if header.Size > maxImageSize {
			utilities.HandleError(w, http.StatusBadRequest, "Image max size is 20MB")
			return
		}

		// Layer 3: io.LimitReader — trustworthy precise enforcement on actual bytes
		// reads up to maxFileSize+1 to detect if file exceeds the limit
		limitedReader := io.LimitReader(file, maxImageSize+1)
		fileBytes, err := io.ReadAll(limitedReader)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if int64(len(fileBytes)) > maxImageSize {
			utilities.HandleError(w, http.StatusRequestEntityTooLarge, "Image max size is 20MB")
			return
		}

		// Check if file is valid image
		buffer := make([]byte, 512)
		file.Seek(0, 0) // without it, Read may give EOF error
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not save image")
			return
		}
		// Reset file pointer so it can be read again later
		if _, err := file.Seek(0, 0); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not save image")
			return
		}
		contentType := http.DetectContentType(buffer)
		if !strings.HasPrefix(contentType, "image/") { // svg not handled: complicated + unsafe xml
			utilities.HandleError(w, http.StatusBadRequest, "Invalid image type")
			return
		}

		imageUri, err = utilities.SaveImage(file, header)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not save image")
			return
		}
	}

	tx, err := database.Database.Begin()
	if err != nil {
		utilities.HandleError(w, http.StatusInternalServerError, "Could not create post")
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CreateComment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/comments/create" {
		utilities.HandleError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodPost {
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postId := strings.TrimSpace(r.FormValue("postId"))
	text := strings.TrimSpace(r.FormValue("text"))

	if text == "" {
		utilities.HandleError(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	if postId == "" {
		utilities.HandleError(w, http.StatusBadRequest, "Invalid post")
		return
	}
	if len(text) > 1000 {
		utilities.HandleError(w, http.StatusBadRequest, "Comment cannot exceed 1000 characters")
		return
	}

	_, err := strconv.Atoi(postId)
	if err != nil {
		utilities.HandleError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	cookie, _ := r.Cookie("session_id")
	userId, err := models.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}
	if _, err = database.Database.Exec(
		"INSERT INTO comments (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)",
		userId,
		postId,
		time.Now(),
		text,
	); err != nil {
		utilities.HandleError(w, http.StatusInternalServerError, "Could not create comment")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PostResolver
func PostResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	userId, err := models.GetUserIDFromCookie(cookie.Value)
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
		if err := models.ReactToPost(userId, postId, 1); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not react to post")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "dislike":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := models.ReactToPost(userId, postId, -1); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not react to post")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "delete":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := models.DeletePost(postId, userId); err != nil {
			utilities.HandleError(w, http.StatusForbidden, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		utilities.HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}
}

// CommentResolver
func CommentResolver(w http.ResponseWriter, r *http.Request) {
	endpoint := r.PathValue("endpoint")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	userId, err := models.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.HandleError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	commentId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.HandleError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	switch endpoint {
	case "like":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := models.ReactToComment(userId, commentId, 1); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not react to comment")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "dislike":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := models.ReactToComment(userId, commentId, -1); err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Could not react to comment")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "delete":
		if r.Method != http.MethodPost {
			utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		if err := models.DeleteComment(commentId, userId); err != nil {
			utilities.HandleError(w, http.StatusForbidden, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		utilities.HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}
}
