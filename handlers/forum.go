package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"forum/database"
	api "forum/forum-api"
)

type TemplateData struct {
	IsLoggedIn bool
	User       User
	Posts      []api.Post
}

func Forum(w http.ResponseWriter, r *http.Request) {
	// Validate route
	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse template
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Template error")
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		// ParseForm parses the raw query from the URL and updates r.Form
		HandleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	categories := r.Form["categories"]
	isLiked := r.FormValue("my-liked-post") == "true"  // 1. needs format
	isByMe := r.FormValue("my-creat-postes") == "true" // 2. needs format

	var (
		user   User
		userId int
	)

	// Try to get user from cookie
	cookie, err := r.Cookie("session_id")
	if err == nil {
		userId, _ = api.GetUserIDFromCookie(cookie.Value)

		user, err = getUser(cookie.Value)
		if err != nil {
			// log.Println("error getting user:", err)
			user = User{}
		}
	}

	// Fetch posts
	posts, err := api.GetFilteredPosts(userId, categories, isLiked, isByMe)
	if err != nil {
		log.Println("error getting posts:", err)
		posts = []api.Post{} // fallback
	}

	// Mark liked posts
	if user.Id != 0 {
		api.CheckLikedPosts(posts, user.Id)
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
		HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func getUser(sessionId string) (User, error) {
	var user User
	err := database.Database.QueryRow(
		"SELECT u.id, u.name, u.email FROM USERS u INNER JOIN SESSIONS s ON s.user_id = u.id WHERE s.id = ? AND s.expires_at > DATETIME('now')",
		sessionId,
	).Scan(&user.Id, &user.Name, &user.Email) //, &user.Password)
	// + reading password problem! -> fortunately not needed
	// use sql.NullString or *string if needed
	return user, err
}
