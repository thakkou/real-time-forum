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

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		utilities.WriteJSON(w, 404, `path not found`, nil)
		return
	}
	if r.Method != "POST" {
		utilities.WriteJSON(w, 405, `method not allowed`, nil)
		return
	}
	user := models.User{}
	identifier := strings.TrimSpace(r.FormValue("email")) // email or username
	password := r.FormValue("password")

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

}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		utilities.HandleError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodPost {
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil { // http.ErrNoCookie
		return
	}

	err = utilities.DeleteSession(cookie.Value)
	// + need to remove cookie from storage
	if err != nil {
		log.Println(err)
	}

	http.SetCookie(w, &http.Cookie{ // delete cookie ------------------- TODO: function already exists
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther) // or to login
}

const RULES string = `
. name (valid)  : 2 ~ 50  chars
. age           : 1 <= x <= 99
. email (valid) : 5 ~ 100 chars
. password      : 6 ~ 20  chars`

// Register
func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		utilities.WriteJSON(w, 404, `path not found`, nil)
		return
	}

	if r.Method != "POST" {
		utilities.WriteJSON(w, 405, `method not allowed`, nil)
		return
	}
	age, errAge := strconv.Atoi(strings.TrimSpace(r.FormValue("age")))
	if errAge != nil {
		utilities.WriteJSON(w, 400, "invalid age", nil)
		return
	}

	user := models.User{
		Nickname:        strings.TrimSpace(r.FormValue("nickname")),
		FirstName:       strings.TrimSpace(r.FormValue("first_name")),
		LastName:        strings.TrimSpace(r.FormValue("last_name")),
		Gender:          strings.TrimSpace(r.FormValue("gender")),
		Age:             age,
		Email:           strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}

	fmt.Println("validat empty")
	// ** check emptiness
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
			utilities.WriteJSON(
				w,
				http.StatusBadRequest,
				field.Name+" is required",
				nil,
			)
			return
		}
	}
	fmt.Println("start validate")

	if !utilities.IsValidName(user.Nickname) ||
		!utilities.IsValidName(user.FirstName) ||
		!utilities.IsValidName(user.LastName) ||
		!utilities.IsValidAge(user.Age) ||
		!utilities.IsValidGender(user.Gender) ||
		!utilities.IsValidEmail(user.Email) ||
		!utilities.IsValidPassword(user.Password) {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid age", RULES)
		return
	}

	if user.Password != user.ConfirmPassword {
		utilities.WriteJSON(w, http.StatusBadRequest, "Password not confirmed", RULES)
		return
	}

	// 3. Check email availability
	var emailExists bool
	err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email,
	).Scan(&emailExists)

	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}
	if emailExists {
		utilities.WriteJSON(w, http.StatusBadRequest, "Email already exist", nil)
		return
	}

	var nicknameExists bool
	err = database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE nickname = ? COLLATE NOCASE)", user.Nickname,
	).Scan(&nicknameExists)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error,nickname", nil)
		return
	}
	if nicknameExists {
		utilities.WriteJSON(w, http.StatusBadRequest, "Username already taken", nil)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error,hashing", nil)
		return
	}

	// 7. Create user
	_, err = database.Database.Exec(
		"INSERT INTO users (nickname, firstname, lastname, age, gender, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Nickname,
		user.FirstName,
		user.LastName,
		user.Age,
		user.Gender,
		user.Email,
		string(hashedPassword),
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error,insertion", err)
		return
	}

	utilities.WriteJSON(w, 200, "registration sucess", user)
}
