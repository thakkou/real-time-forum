package models

import "forum/database"

type User struct {
	Id       int    `json:"id"`   // check for google
	Name     string `json:"name"` // name or username: problem for providers!
	Email    string `json:"email"`
	Password string

	// Not Used:
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`

	// Not Stored:
	ConfirmPassword string
	Message         string
	// Picture string `json:"picture"`    // gmail picture: sometimes cannot be loaded!
	// Avatar  string `json:"avatar_url"` // github avatar
}

// GetUser
func GetUser(sessionId string) (User, error) {
	var user User
	err := database.Database.QueryRow(
		"SELECT u.id, u.name, u.email FROM USERS u INNER JOIN SESSIONS s ON s.user_id = u.id WHERE s.id = ? AND s.expires_at > DATETIME('now')",
		sessionId,
	).Scan(&user.Id, &user.Name, &user.Email) //, &user.Password)
	// + reading password problem! -> fortunately not needed
	// use sql.NullString or *string if needed
	return user, err
}
