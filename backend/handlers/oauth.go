package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"forum/database"
	"forum/utilities"

	"github.com/google/uuid"
)

var (
	// Google
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	redirectURI          = "http://localhost:8080/auth/google/callback"
	// Github
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
)

// OAuthLogin
func OAuthLogin(w http.ResponseWriter, r *http.Request) {
	var baseURL, client_id, scope string

	provider := r.PathValue("provider")
	switch provider {
	case "google":
		baseURL = "https://accounts.google.com/o/oauth2/v2/auth"
		client_id = GOOGLE_CLIENT_ID
		scope = "openid email profile"

	case "github":
		baseURL = "https://github.com/login/oauth/authorize"
		client_id = GITHUB_CLIENT_ID
		scope = "read:user user:email"

	default:
		utilities.WriteJSON(w, http.StatusNotFound, "Unknown endpoint", nil)
	}

	// Generate a random state token to prevent CSRF
	// state := generateState()
	// stateStore[state] = true

	params := url.Values{}
	params.Add("client_id", client_id)

	// only for google ?!
	if provider == "google" {
		params.Add("redirect_uri", redirectURI)
	}

	params.Add("response_type", "code")
	params.Add("scope", scope)

	// only for google ?!
	params.Add("access_type", "offline")

	// params.Add("prompt", "consent")
	// Forces account chooser to appear every time (optional)
	params.Set("prompt", "select_account")

	// params.Set("state", state)

	authURL := baseURL + "?" + params.Encode()

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// OAuthCallback
func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	var tokenURL, client_id, client_secret, userInfoURL string

	provider := r.PathValue("provider")
	switch provider {
	case "google":
		tokenURL = "https://oauth2.googleapis.com/token"
		client_id = GOOGLE_CLIENT_ID
		client_secret = GOOGLE_CLIENT_SECRET
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

	case "github":
		tokenURL = "https://github.com/login/oauth/access_token"
		client_id = GITHUB_CLIENT_ID
		client_secret = GITHUB_CLIENT_SECRET
		userInfoURL = "https://api.github.com/user"

	default:
		utilities.WriteJSON(w, http.StatusNotFound, "Unknown endpoint", nil)
	}

	// 1. Validate state to prevent CSRF
	// state := r.URL.Query().Get("state")
	// if !stateStore[state] {
	// 	http.Error(w, "invalid state", http.StatusBadRequest)
	// 	return
	// }
	// delete(stateStore, state) // one-time use

	// 2. Check for errors (e.g. user denied consent)
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		http.Error(w, "auth error: "+errMsg, http.StatusUnauthorized)
		return
	}

	// 3. Exchange authorization code for tokens
	code := r.URL.Query().Get("code")
	tokenData, err := utilities.ExchangeCode(provider, tokenURL, client_id, client_secret, code)
	if err != nil {
		http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	accessToken := tokenData.AccessToken

	// 4. Fetch user info
	user, err := utilities.FetchUserInfo(userInfoURL, accessToken)
	if provider == "github" {
		user.Email, _ = utilities.FetchGithubUserEmail(accessToken)
	}
	if err != nil {
		http.Error(w, "failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Processing claimed user info
	if provider == "github" {
		user.FirstName, user.LastName, _ = strings.Cut(user.Name, " ")
		user.Name = user.Login
	}
	user.Age = 18                            // default
	user.Gender = "male"                     // default
	user.Email = strings.ToLower(user.Email) // case insensitivity

	// 6. check email and username availability
	// Similar to register.go
	var emailExists bool
	var nameExists bool = true
	err = database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email,
	).Scan(&emailExists)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Database error", nil)
		return
	}
	var newUsername string
	var i int
	if !emailExists {
		// check username & generate alternatives
		for nameExists {
			// add suffix for username, and result shouldnt exist
			if newUsername == "" {
				newUsername = user.Name
			} else {
				newUsername = user.Name + "+" + strconv.Itoa(i)
				i++
			}
			err = database.Database.QueryRow(
				"SELECT EXISTS(SELECT 1 FROM users WHERE name = ? COLLATE NOCASE)", newUsername,
			).Scan(&nameExists)
			if err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "Database error", nil)
				return
			}
		}

		// create user
		user.Name = newUsername
		var gender_id int
		if user.Gender == "male" {
			gender_id = 1
		}
		_, err = database.Database.Exec(
			"INSERT INTO users (name, firstname, lastname, age, gender, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)",
			user.Name,
			user.FirstName,
			user.LastName,
			user.Age,
			gender_id,
			user.Email,
			nil,
		)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "Server error", nil)
			return
		}
	}

	// 7. generate toke based on email not username (username is the one in db)
	// ************* code from login.go ****************
	var userID int

	err = database.Database.QueryRow(
		"SELECT id FROM users WHERE email = ?", user.Email,
	).Scan(&userID)
	// if err != nil {
	// 	user := User{Message: "Invalid email/username or password"}
	// 	RenderTemplate(w, 400, "login.html", user)
	// 	return
	// }

	// Delete any existing sessions for this user
	_, err = database.Database.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Server error", nil)
		return
	}

	sessionID := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour)

	_, err = database.Database.Exec(
		"INSERT INTO SESSIONS (id, expires_at, user_id) VALUES (?, ?, ?)",
		sessionID, expiration, userID,
	)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "Server error", nil)
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
	// **************************************************
}
