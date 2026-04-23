package utilities

import (
	"net/mail"
	"regexp"
	"strings"
)

func IsValidName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_ ]{1,49}$`)
	// disallowing multiple spaces
	return re.MatchString(name) && !strings.Contains(name, "  ")
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return len(email) >= 5 && len(email) <= 100 && (err == nil)
}

func IsValidPassword(password string) bool {
	return len(password) >= 6 && len(password) <= 20
}

func IsValidAge(age int) bool {
	return age >= 1 && age <= 99
}

func IsValidGender(gender string) bool {
	return gender == "male" || gender == "female"
}
