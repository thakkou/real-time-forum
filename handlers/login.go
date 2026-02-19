package zone

import (
	"database/sql"
	"html/template"
	"net/http"
	"zone/database"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Login" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")

		var dbPassword string

		err := database.Database.QueryRow(
			"SELECT password FROM users WHERE email = ?",
			email,
		).Scan(&dbPassword)

		if err == sql.ErrNoRows {
			http.Error(w, "User not found", 401)
			return
		}

		if err != nil {
			http.Error(w, "Server error", 500)
			return
		}

		if password != dbPassword {
			http.Error(w, "Wrong password", 401)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/Login.html")
	if err != nil {
		http.Error(w, "Template error", 500)
		return
	}

	tmpl.Execute(w, nil)
}
