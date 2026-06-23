package utilities

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
)

// IsValidName
func IsValidName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_ ]{1,49}$`)
	return re.MatchString(name) && !strings.Contains(name, "  ")
}

// IsValidEmail
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return len(email) >= 5 && len(email) <= 100 && (err == nil)
}

// IsValidPassword
func IsValidPassword(password string) bool {
	return len(password) >= 6 && len(password) <= 20
}

func IsValidAge(age any) bool {
	var a int

	switch v := age.(type) {
	case int:
		a = v

	case int64:
		a = int(v)

	case float64:
		a = int(v)

	case string:
		// try convert string → int
		parsed, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return false
		}
		a = parsed

	default:
		return false
	}

	return a >= 1 && a <= 99
}

// IsValidGender
func IsValidGender(gender string) bool {
	gender = strings.ToLower(gender)
	return gender == "male" || gender == "female"
}

func ToInt(v any) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int64:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}
