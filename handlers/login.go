package handlers

import (
	"net/http"
	"strings"
	"time"

	"forum/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, 200, "login.html", nil)

	case http.MethodPost:
		identifier := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		// Basic input validation
		if identifier == "" || password == "" {
			HandleError(w, http.StatusBadRequest, "Email/username and password are required")
			return
		}

		var userID int
		var hashedPassword string

		err := database.Database.QueryRow(
			"SELECT id, password FROM users WHERE email = ? OR name = ?", identifier, identifier,
		).Scan(&userID, &hashedPassword)
		if err != nil {
			user := User{Message: "Invalid email/username or password"}
			RenderTemplate(w, 400, "login.html", user)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			user := User{Message: "Invalid email/username or password"}
			RenderTemplate(w, 401, "login.html", user)
			return
		}

		// Delete any existing sessions for this user
		_, err = database.Database.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Server error")
			return
		}

		sessionID := uuid.New().String()
		expiration := time.Now().Add(24 * time.Hour)

		_, err = database.Database.Exec(
			"INSERT INTO SESSIONS (id, expires_at, user_id) VALUES (?, ?, ?)",
			sessionID, expiration, userID,
		)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Server error")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Expires:  expiration,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}