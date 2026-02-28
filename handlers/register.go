package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"forum/database"

	"golang.org/x/crypto/bcrypt"
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

	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("templates/register.html")
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Template error")
			return
		}
		var buf bytes.Buffer
		if err := t.Execute(&buf, nil); err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)
	case http.MethodPost:
		user := User{
			Name:     strings.TrimSpace(r.FormValue("name")),
			Email:    strings.TrimSpace(r.FormValue("email")),
			Password: r.FormValue("password"),
		}

		// Input validation
		if user.Name == "" || user.Email == "" || user.Password == "" {
			HandleError(w, http.StatusBadRequest, "All fields are required")
			return
		}
		if len(user.Name) < 2 || len(user.Name) > 50 {
			HandleError(w, http.StatusBadRequest, "Name must be between 2 and 50 characters")
			return
		}
		if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
			HandleError(w, http.StatusBadRequest, "Invalid email address")
			return
		}
		if len(user.Password) < 6 || len(user.Password) > 20 {
			HandleError(w, http.StatusBadRequest, "Password must be between 6 and 20 characters")
			return
		}

		// Check if email already exists
		var exists bool
		err := database.Database.QueryRow(
			"SELECT EXISTS(SELECT * FROM users WHERE email = ?)", user.Email,
		).Scan(&exists)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Database error")
			return
		}
		if exists {
			HandleError(w, http.StatusConflict, "Email already taken")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Password hashing error")
			return
		}

		_, err = database.Database.Exec(
			"INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
			user.Name,
			user.Email,
			string(hashedPassword),
		)
		// create session if you want to redirect to its page
		if err != nil {
			// log.Println(err.Error())
			HandleError(w, http.StatusInternalServerError, "Could not create account")
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
