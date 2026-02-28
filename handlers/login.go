package handlers

import (
	"bytes"
	"html/template"
	"log"
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
		t, err := template.ParseFiles("templates/login.html")
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
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		// Basic input validation
		if email == "" || password == "" {
			HandleError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		var userID int
		var hashedPassword string

		err := database.Database.QueryRow(
			"SELECT id, password FROM users WHERE email = ?", email,
		).Scan(&userID, &hashedPassword)
		if err != nil { // sql.ErrNoRows
			// Don't reveal whether email exists or not
			HandleError(w, http.StatusUnauthorized, "Invalid email or password") // 401 ?
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			HandleError(w, http.StatusUnauthorized, "Invalid email or password") // 401
			return
		}

		/////////////////////////
		// Need a mechanism to remove expired sessions from database (from time to time ?!)
		/////////////////////////

		sessionID := uuid.New().String()             // OR: uuid.NewString() // unique ?
		expiration := time.Now().Add(24 * time.Hour) // DATETIME('now', '+24 hours')

		_, err = database.Database.Exec(
			"INSERT INTO SESSIONS (id, expires_at, user_id) VALUES (?, ?, ?)",
			sessionID, expiration, userID,
		)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Server error") // message
			log.Println(err)
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
