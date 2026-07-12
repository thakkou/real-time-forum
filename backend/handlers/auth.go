package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/database"
	"forum/models"
	"forum/utilities"
	"forum/ws"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/login" {
		utilities.WriteJSON(w, http.StatusNotFound, "path not found", nil)
		return
	}

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		utilities.WriteJSON(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
		return
	}

	type LoginModel struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	userLog, err := utilities.ReadJSONRequest[LoginModel](r)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if userLog.Identifier == "" || userLog.Password == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "bad credentials", nil)
		return
	}

	var (
		userID         int
		nickname       string
		hashedPassword sql.NullString
	)

	err = database.Database.QueryRow(
		`SELECT id, nickname, password
		 FROM users
		 WHERE email = ? OR nickname = ?`,
		userLog.Identifier,
		userLog.Identifier,
	).Scan(&userID, &nickname, &hashedPassword)
	if err != nil {
		log.Printf("[LOGIN] User not found: %q (%v)", userLog.Identifier, err)

		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid email/username or password.", nil)
		return
	}

	if !hashedPassword.Valid {
		log.Printf("[LOGIN] User %q has no valid password hash", userLog.Identifier)

		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid email/username or password.", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword.String),
		[]byte(userLog.Password),
	); err != nil {
		log.Printf("[LOGIN] Invalid password for user %q password=%q", userLog.Identifier, userLog.Password)
		utilities.WriteJSON(w, http.StatusUnauthorized, "Invalid email/username or password.", nil)
		return
	}

	// Remove old sessions
	_, err = database.Database.Exec(
		"DELETE FROM sessions WHERE user_id = ?",
		userID,
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	// Create new session
	sessionID := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour)

	_, err = database.Database.Exec(
		"INSERT INTO sessions (id, expires_at, user_id) VALUES (?, ?, ?)",
		sessionID,
		expiration,
		userID,
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}
	ws.NotifyUser(strconv.Itoa(userID), "force_logout", nil)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiration,
	})

	utilities.WriteJSON(w, http.StatusOK, "Login Success", map[string]any{
		"user_id":  userID,
		"nickname": nickname,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/logout" {
		utilities.WriteJSON(w, 404, `path not found`, nil)
		return
	}
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, `Method not allowed`, nil)

		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil { // http.ErrNoCookie
		return
	}

	err = utilities.DeleteSession(cookie.Value)
	if err != nil {
		log.Println(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	utilities.WriteJSON(w, 201, `log out succes`, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Check route
	if r.URL.Path != "/api/register" {
		utilities.WriteJSON(w, http.StatusNotFound, "path not found", nil)
		return
	}

	// Check method
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	// Content type (optional but fine to keep)
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		utilities.WriteJSON(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
		return
	}
	// ✅ REPLACED PART (clean)
	user, err := utilities.ReadJSONRequest[models.User](r)
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// Normalize input
	user.Nickname = strings.TrimSpace(user.Nickname)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Gender = strings.TrimSpace(user.Gender)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	// Check required fields
	fields := []struct {
		Name  string
		Value string
	}{
		{"nickname", user.Nickname},
		{"first_name", user.FirstName},
		{"last_name", user.LastName},
		{"gender", user.Gender},
		{"email", user.Email},
		{"password", user.Password},
		{"confirm_password", user.ConfirmPassword},
	}

	for _, field := range fields {
		if field.Value == "" {
			utilities.WriteJSON(w, http.StatusBadRequest, field.Name+" is required", nil)
			return
		}
	}
	fmt.Println("pass", user.Password, user.Password == user.ConfirmPassword)
	// Validate fields
	// Validate fields one by one
	if !utilities.IsValidName(user.Nickname) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid nickname: use only letters (2-50 characters)", nil)
		return
	}

	if !utilities.IsValidName(user.FirstName) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid first name: use only letters (2-50 characters)", nil)
		return
	}

	if !utilities.IsValidName(user.LastName) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid last name: use only letters (2-50 characters)", nil)
		return
	}

	if !utilities.IsValidAge(user.Age) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid age: must be between 1 and 99", nil)
		return
	}

	if !utilities.IsValidGender(user.Gender) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid gender: must be 'Male' or 'Female'", nil)
		return
	}

	if !utilities.IsValidEmail(user.Email) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid email: must be a valid email address", nil)
		return
	}

	if !utilities.IsValidPassword(user.Password) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid password: must be 6-25 characters", nil)
		return
	}

	// Confirm password
	if user.Password != user.ConfirmPassword {
		utilities.WriteJSON(w, http.StatusBadRequest, "password not confirmed", nil)
		return
	}

	// Check email exists
	var emailExists bool
	err = database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		user.Email,
	).Scan(&emailExists)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	if emailExists {
		utilities.WriteJSON(w, http.StatusBadRequest, "email already exists", nil)
		return
	}

	// Check nickname exists
	var nicknameExists bool
	err = database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE nickname = ? COLLATE NOCASE)",
		user.Nickname,
	).Scan(&nicknameExists)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	if nicknameExists {
		utilities.WriteJSON(w, http.StatusBadRequest, "username already taken", nil)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	// Insert user
	_, err = database.Database.Exec(
		`INSERT INTO users (nickname, firstname, lastname, age, gender, email, password)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Nickname,
		user.FirstName,
		user.LastName,
		user.Age,
		user.Gender,
		user.Email,
		string(hashedPassword),
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}
	type RegisterResponse struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}
	response := RegisterResponse{
		Nickname: user.Nickname,
		Email:    user.Email,
	}

	utilities.WriteJSON(w, http.StatusOK, "registration success", response)
}
