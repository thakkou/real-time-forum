package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
)

type User struct {
	Name     string
	Email    string
	Password string
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method == http.MethodPost {

		user := User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		_, err := database.Database.Exec(
			"INSERT INTO users (userName, email, password) VALUES (?, ?, ?)",
			user.Name,
			user.Email,
			user.Password,
		)
		if err != nil {
			HandleError(w, 500, "Internal Server Error")
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		HandleError(w, 500, "Template error")
		return
	}

	tmpl.Execute(w, nil)
}
