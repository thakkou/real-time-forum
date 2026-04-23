package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"forum/models"
	"forum/utilities"
)

// TODO: moving to utilities caused cyclic import
type TemplateData struct {
	IsLoggedIn bool
	User       models.User
	Posts      []models.Post
}

// Forum
func Forum(w http.ResponseWriter, r *http.Request) {
	// Validate route
	if r.URL.Path != "/" {
		utilities.HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodGet {
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse template
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		utilities.HandleError(w, http.StatusInternalServerError, "Template error")
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		// ParseForm parses the raw query from the URL and updates r.Form
		utilities.HandleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	categories := r.Form["categories"]
	isLiked := r.FormValue("my-liked-post") == "true"  // 1. needs format
	isByMe := r.FormValue("my-creat-postes") == "true" // 2. needs format

	var (
		user   models.User
		userId int
	)

	// Try to get user from cookie
	cookie, err := r.Cookie("session_id")
	if err == nil {
		userId, _ = models.GetUserIDFromCookie(cookie.Value)

		user, err = models.GetUser(cookie.Value)
		if err != nil {
			// log.Println("error getting user:", err)
			user = models.User{}
		}
	}

	// Fetch posts
	posts, err := models.GetFilteredPosts(userId, categories, isLiked, isByMe)
	if err != nil {
		log.Println("error getting posts:", err)
		posts = []models.Post{} // fallback
	}

	// Mark liked posts
	if user.Id != 0 {
		models.CheckLikedPosts(posts, user.Id)
	}

	// Prepare template data
	data := TemplateData{
		IsLoggedIn: user.Id != 0,
		User:       user,
		Posts:      posts,
	}

	// Render template
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Println("template execute error:", err)
		utilities.HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
