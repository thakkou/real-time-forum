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
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
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
	text = strings.ReplaceAll(text, "\r\n", "\n")
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
	file, header, err := r.FormFile("image")
	if err != nil {
		// post without image is handled!
		fmt.Println("No image uploaded, continuing without it")
	} else {
		defer file.Close()

		// Layer 2: header.Size — fast pre-check, avoids reading the file at all (depends on content-length)
		// not 100% trustworthy (client-declared) but useful to reject obviously large files early
		if header.Size > maxImageSize {
			HandleError(w, http.StatusBadRequest, "Image max size is 20MB")
			return
		}

		// Layer 3: io.LimitReader — trustworthy precise enforcement on actual bytes
		// reads up to maxFileSize+1 to detect if file exceeds the limit
		limitedReader := io.LimitReader(file, maxImageSize+1)
		fileBytes, err := io.ReadAll(limitedReader)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if int64(len(fileBytes)) > maxImageSize {
			HandleError(w, http.StatusRequestEntityTooLarge, "Image max size is 20MB")
			return
		}

		// Check if file is valid image
		buffer := make([]byte, 512)
		file.Seek(0, 0) // without it, Read may give EOF error
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			HandleError(w, http.StatusInternalServerError, "Could not save image")
			return
		}
		// Reset file pointer so it can be read again later
		if _, err := file.Seek(0, 0); err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not save image")
			return
		}
		contentType := http.DetectContentType(buffer)
		if !strings.HasPrefix(contentType, "image/") { // svg not handled: complicated + unsafe xml
			HandleError(w, http.StatusBadRequest, "Invalid image type")
			return
		}

		imageUri, err = SaveImage(file, header)
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

func SaveImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Ensure uploads directory exists
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("unable to create uploads directory: %w", err)
	}

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
