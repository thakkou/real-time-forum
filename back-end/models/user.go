package models

import (
	"forum/database"
)

type User struct {
	Id int
	// register
	Password        string `json:"password"`
	Login           string `json:"login"`
	Name            string `json:"name"`
	ConfirmPassword string `json:"confirm_password"`
	Nickname        string `json:"nickname"`
	// github + (google requires other apis => 'name')
	FirstName string `json:"first_name"` // google
	LastName  string `json:"last_name"`  // google
	Age       any    `json:"age"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`

	Message string // (NOT STORED)
	// google oauth

	// Picture string `json:"picture"`    // gmail picture: sometimes cannot be loaded!
	// Avatar  string `json:"avatar_url"` // github avatar
}

// GetUser
func GetUser(sessionId string) (User, error) {
	var user User
	err := database.Database.QueryRow(
		"SELECT u.id, u.name, u.email FROM USERS u INNER JOIN SESSIONS s ON s.user_id = u.id WHERE s.id = ? AND s.expires_at > DATETIME('now')",
		sessionId,
	).Scan(&user.Id, &user.Name, &user.Email)

	return user, err
}
