package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"forum/database"
	"forum/models"
	"forum/utilities"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		utilities.HandleError(w, http.StatusNotFound, "Page Not Found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		utilities.RenderTemplate(w, 200, "login.html", nil)

	case http.MethodPost:
		// 1. Get form data
		user := models.User{}
		identifier := strings.TrimSpace(r.FormValue("email")) // email or username
		password := r.FormValue("password")

		// Basic input validation
		if identifier == "" || password == "" {
			user.Message = "All fields are required."
			utilities.RenderTemplate(w, http.StatusBadRequest, "login.html", user) // 400
			return
		}

		var userID int
		var hashedPassword sql.NullString

		err := database.Database.QueryRow(
			"SELECT id, password FROM users WHERE email = ? OR name = ?", identifier, identifier,
		).Scan(&userID, &hashedPassword)
		if err != nil {
			user.Message = "Invalid email/username or password."
			utilities.RenderTemplate(w, http.StatusBadRequest, "login.html", user) // 400
			return
		}

		if !hashedPassword.Valid {
			user.Message = "Account registred by provider."                          // not good practice
			utilities.RenderTemplate(w, http.StatusUnauthorized, "login.html", user) // 401
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword.String), []byte(password)); err != nil {
			user.Message = "Invalid email/username or password."
			utilities.RenderTemplate(w, http.StatusUnauthorized, "login.html", user) // 401
			return
		}

		// Delete any existing sessions for this user
		_, err = database.Database.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		sessionID := uuid.New().String()
		expiration := time.Now().Add(24 * time.Hour)

		_, err = database.Database.Exec(
			"INSERT INTO SESSIONS (id, expires_at, user_id) VALUES (?, ?, ?)",
			sessionID, expiration, userID,
		)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
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
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}
