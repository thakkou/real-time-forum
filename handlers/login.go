package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"forum/database"

	"github.com/google/uuid"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var dbPassword, userName string
		err := database.Database.QueryRow(
			"SELECT password, userName FROM users WHERE email = ?",
			email,
		).Scan(&dbPassword, &userName)
		if err == sql.ErrNoRows {
			HandleError(w, 401, "User not found")
			return
		}
		if err != nil {
			HandleError(w, 500, "User not found")
			return
		}

		if password != dbPassword {
			HandleError(w, 401, "Wrong password")

			return
		}

		sessionID := uuid.NewString()
		_, err = database.Database.Exec(
			"UPDATE users SET session = ?, dateexpired = DATETIME('now', '+24 hours') WHERE email = ?",
			sessionID, email,
		)
		if err != nil {
			HandleError(w, 500, "Server error")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		HandleError(w, 500, "Template error")
		return
	}

	tmpl.Execute(w, nil)
}
