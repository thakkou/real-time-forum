package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum/database"
	api "forum/forum-api"
	"forum/helper"
)

// func HandleStatic(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path == "/static" || r.URL.Path == "/static/" {
// 		HandleError(w, http.StatusNotFound, "Not Found")
// 		return
// 	}

// 	filePath := filepath.Join("static", r.URL.Path[len("/static/"):])

// 	// Prevent path traversal (e.g. /static/../../secret)
// 	cleanPath := filepath.Clean(filePath)
// 	if len(cleanPath) < len("static") || cleanPath[:len("static")] != "static" {
// 		HandleError(w, http.StatusForbidden, "Forbidden")
// 		return
// 	}

// 	if _, err := os.Stat(filePath); err != nil {
// 		HandleError(w, http.StatusNotFound, "Not Found")
// 		return
// 	}

// 	http.ServeFile(w, r, filePath)
// }

// =======================================================================

type TemplateData struct {
	IsLoggedIn bool
	User       User
	Posts      []api.Post
}

func Forum(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // root dir or something else
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
	r.ParseForm() // must call first
	// get posts
	categories := r.Form["categories"]
	isLiked := r.FormValue("my-liked-post")
	isByMe := r.FormValue("my-creat-postes")
	var posts []api.Post
	if len(categories) == 0 && isLiked != "true" && isByMe != "true" {
		posts, _ = api.GetPosts()
	} else {
		//get the user ID
				cookie, err := r.Cookie("session_id")
				if(err != nil){
					fmt.Println("No session cookie found, treating as not logged in")
							posts, _ = api.GetPosts()

					
				}
userId,_:=helper.GetUserIDFromCookie(cookie.Value)	
	

		posts, _ = api.GetFiltrtPOst(userId, categories, isLiked == "true", isByMe == "true")
	}

	var buf bytes.Buffer
	cookie, err := r.Cookie("session_id")
	if err != nil { // http.ErrNoCookie
		if err := t.Execute(&buf, TemplateData{Posts: posts}); err != nil {
			log.Printf("home template execute error: %v", err)
			return // ?
		}
		// nil works fine without using: TemplateData{}
		buf.WriteTo(w)
		return
	}

	// get User
	user, err := getUser(cookie.Value)
	if err != nil { // sql.ErrNoRows
		// what is the default behavior when session cookie not found -> serve as not logged in ?
		if err := t.Execute(&buf, TemplateData{Posts: posts}); err != nil {
			log.Printf("home template execute error: %v", err)
			return // ?
		}
		buf.WriteTo(w)
		return
	}

	// check each post if liked or disliked bu the current user
	api.CheckLikedPosts(posts, user.Id)

	data := TemplateData{
		IsLoggedIn: true,
		User:       user,
		Posts:      posts,
	}
	err = t.Execute(&buf, data)
	if err != nil {
		log.Println(err)
		HandleError(w, http.StatusInternalServerError, "Internal server error")
		// send err.Error() as message !
		return
	}

	// send successful response
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
