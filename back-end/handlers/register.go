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

const RULES string = `
. name (valid)  : 2 ~ 50  chars
. age           : 1 <= x <= 99
. email (valid) : 5 ~ 100 chars
. password      : 6 ~ 20  chars`

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
		// 1. Get registration form data
		user := models.User{
			Name:            strings.TrimSpace(r.FormValue("name")), // nickname
			FirstName:       strings.TrimSpace(r.FormValue("firstname")),
			LastName:        strings.TrimSpace(r.FormValue("lastname")),
			Gender:          strings.TrimSpace(r.FormValue("gender")),
			Email:           strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),
		}
		var err error
		user.Age, err = strconv.Atoi(r.FormValue("age"))
		if err != nil {
			user.Message = "Age: Not a number"
			utilities.RenderTemplate(w, http.StatusBadRequest, "register.html", user) // 400
			return
		}

		// 2. Sanitize form data
		// ** check emptiness
		for _, f := range []string{user.Name, user.FirstName, user.LastName, user.Gender, user.Email, user.Password, user.ConfirmPassword} {
			if f == "" {
				user.Message = "All fields are required"
				utilities.RenderTemplate(w, http.StatusBadRequest, "register.html", user) // 400
				return
			}
		}
		// ** check validity
		if !utilities.IsValidName(user.Name) ||
			!utilities.IsValidName(user.FirstName) ||
			!utilities.IsValidName(user.LastName) ||
			!utilities.IsValidAge(user.Age) ||
			!utilities.IsValidGender(user.Gender) ||
			!utilities.IsValidEmail(user.Email) ||
			!utilities.IsValidPassword(user.Password) {
			user.Message = RULES
			utilities.RenderTemplate(w, http.StatusBadRequest, "register.html", user) // 400
			return
		}
		// ** check password match
		if user.Password != user.ConfirmPassword {
			user.Message = "Password not confirmed"
			utilities.RenderTemplate(w, http.StatusBadRequest, "register.html", user) // 400
			return
		}

		// 3. Check email availability
		var emailExists bool
		err = database.Database.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email,
		).Scan(&emailExists)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if emailExists {
			user.Message = "Email already registered"                                 // not good practice
			utilities.RenderTemplate(w, http.StatusBadRequest, "register.html", user) // 400
			return
		}

		// 4. Check username availability
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
			utilities.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if nameExists {
			user.Message = "Username already taken"
			utilities.RenderTemplate(w, http.StatusBadRequest, "register.html", user) // 400
			return
		}

		// 5. Generate password hash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// 6. Gender id
		var gender_id int
		if user.Gender == "male" {
			gender_id = 1
		}

		// 7. Create user
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
		if err != nil {
			utilities.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}
