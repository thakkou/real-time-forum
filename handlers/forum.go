package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"forum/database"
	api "forum/forum-api"
	"forum/helper"
)

type TemplateData struct {
	IsLoggedIn bool
	User User
	Posts []api.Post
}

func Forum(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Template error")
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, http.StatusBadRequest, "Bad request")
		return
	}

	categories := r.Form["categories"]
	isLiked := r.FormValue("my-liked-post")
	isByMe := r.FormValue("my-creat-postes")

	var posts []api.Post

	if len(categories) == 0 && isLiked != "true" && isByMe != "true" {
		posts, err = api.GetPosts()
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Could not load posts")
			return
		}
	} else {
		cookie, cookieErr := r.Cookie("session_id")
		if cookieErr != nil {
			// no session — just show all posts, filters ignored
			posts, err = api.GetPosts()
			if err != nil {
				HandleError(w, http.StatusInternalServerError, "Could not load posts")
				return
			}
		} else {
			userId, err := helper.GetUserIDFromCookie(cookie.Value)
			if err != nil {
				// invalid session — show all posts
				posts, err = api.GetPosts()
				if err != nil {
					HandleError(w, http.StatusInternalServerError, "Could not load posts")
					return
				}
			} else {
				posts, err = api.GetFiltrtPOst(userId, categories, isLiked == "true", isByMe == "true")
				if err != nil {
					HandleError(w, http.StatusInternalServerError, "Could not load posts")
					return
				}
			}
		}
	}

	var buf bytes.Buffer
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err := t.Execute(&buf, TemplateData{Posts: posts}); err != nil {
			log.Printf("home template execute error: %v", err)
			return
		}
		buf.WriteTo(w)
		return
	}

	user, err := getUser(cookie.Value)
	if err != nil {
		if err := t.Execute(&buf, TemplateData{Posts: posts}); err != nil {
			log.Printf("home template execute error: %v", err)
			return
		}
		buf.WriteTo(w)
		return
	}

	api.CheckLikedPosts(posts, user.Id)

	data := TemplateData{
		IsLoggedIn: true,
		User: user,
		Posts: posts,
	}

	if err = t.Execute(&buf, data); err != nil {
		log.Println(err)
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
		"SELECT u.id, u.name, u.email, u.password FROM USERS u INNER JOIN SESSIONS s ON s.user_id = u.id WHERE s.id = ? AND s.expires_at > DATETIME('now')",
		sessionId,
	).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	return user, err
}