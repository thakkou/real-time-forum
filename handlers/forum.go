package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"forum/database"
)

// Serves the CSS file
func CssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "assets/styles.css")
}

type TemplateData struct {
	IsLoggedIn bool
	User       User
}

func Forum(w http.ResponseWriter, r *http.Request) {
	// 1. check path
	// 2. check method

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleError(w, 500, "Template error")
		return
	}

	var buf bytes.Buffer
	cookie, err := r.Cookie("session_id")
	if err != nil { // http.ErrNoCookie
		// serve dashboard with no connection
		t.Execute(&buf, TemplateData{}) // + check error
		buf.WriteTo(w)

		return
	}

	// var user User
	user, err := getUser(cookie.Value)
	if err != nil { // sql.ErrNoRows
		// what is the default behavior when session cookie not found -> serve as not logged in ?
		t.Execute(&buf, TemplateData{}) // + check error
		buf.WriteTo(w)

		// http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	data := TemplateData{
		IsLoggedIn: true,
		User:       user,
	}
	err = t.Execute(&buf, data)
	if err != nil {
		log.Println("Error executing template.")
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
		"SELECT userName, email, password FROM users WHERE session = ? AND dateexpired > DATETIME('now')",
		sessionId,
	).Scan(&user.Name, &user.Email, &user.Password)
	return user, err
}
