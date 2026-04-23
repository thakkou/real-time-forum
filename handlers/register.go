package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/database"
	"forum/models"
	"forum/utilities"

	"golang.org/x/crypto/bcrypt"
)

// Register
func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		utilities.HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		utilities.RenderTemplate(w, 200, "register.html", nil)

	case http.MethodPost:
		user := models.User{
			Name:            strings.TrimSpace(r.FormValue("name")), // nickname
			Email:           strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),

			// Not used:
			FirstName: strings.TrimSpace(r.FormValue("firstname")),
			LastName:  strings.TrimSpace(r.FormValue("lastname")),
			Gender:    strings.TrimSpace(r.FormValue("gender")),
		}
		var err error
		user.Age, err = strconv.Atoi(r.FormValue("age"))
		if err != nil {
			user.Message = "Age: Not a number"
			utilities.RenderTemplate(w, 400, "register.html", user)
			return
		}

		var rules string = `
. username (valid) : 2 ~ 50  chars 
. email (valid)    : 5 ~ 100 chars
. password         : 6 ~ 20  chars` // is there a newline at first of rules ?!

		// 1. check emptiness
		for _, f := range []string{user.Name, user.FirstName, user.LastName, user.Gender, user.Email, user.Password, user.ConfirmPassword} {
			if f == "" {
				user.Message = "All fields are required"
				utilities.RenderTemplate(w, 400, "register.html", user)
				return
			}
		}
		// 2. check validity
		if !utilities.IsValidName(user.Name) ||
			!utilities.IsValidName(user.FirstName) ||
			!utilities.IsValidName(user.LastName) ||
			!utilities.IsValidAge(user.Age) ||
			!utilities.IsValidGender(user.Gender) ||
			!utilities.IsValidEmail(user.Email) ||
			!utilities.IsValidPassword(user.Password) {
			user.Message = rules
			utilities.RenderTemplate(w, 400, "register.html", user)
			return
		}
		// 3. check password match
		if user.Password != user.ConfirmPassword {
			user.Message = "Password not confirmed"
			utilities.RenderTemplate(w, 400, "register.html", user)
			return
		}

		// username and email fields are case-insensitive
		// emails are lowered then stored, while usernames preserve original casing

		// Check email availability
		var emailExists bool
		err = database.Database.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email,
		).Scan(&emailExists)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Database error")
			return
		}
		if emailExists {
			user.Message = "Email already registered" // not good practice
			utilities.RenderTemplate(w, 400, "register.html", user)
			return
		}

		// Check username availability
		// case insensitivity:
		// - LIKE does treat some special characters as wildcards
		// - LOWER() does affect performance on large databases
		// => used 'COLLATE NOCASE' with a DB INDEX
		// - COLLATE means "how to compare and sort text".
		// - LIMITATION: no unicode support (only ascii)
		var nameExists bool
		err = database.Database.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM users WHERE name = ? COLLATE NOCASE)", user.Name,
		).Scan(&nameExists)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Database error")
			return
		}
		if nameExists {
			user.Message = "Username already taken"
			utilities.RenderTemplate(w, 400, "register.html", user)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Password hashing error")
			return
		}

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
			string(hashedPassword),
		)
		// create session if you want to redirect to its page
		if err != nil {
			// log.Println(err.Error())
			utilities.HandleError(w, http.StatusInternalServerError, "Could not create account")
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
