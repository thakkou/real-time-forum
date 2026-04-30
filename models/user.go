package models

import "forum/database"

type User struct {
	Id int // (NOT USED) 'id' github + 'sub' for google

	Login     string `json:"login"`
	Name      string `json:"name"`        // github + (google requires other apis => 'name')
	FirstName string `json:"given_name"`  // google
	LastName  string `json:"family_name"` // google

	Age    int    // requires other apis: default 18 (+ age changes every year)
	Gender string // requires other apis: default male

	Email    string `json:"email"`
	Password string

	ConfirmPassword string // (NOT STORED)
	Message         string // (NOT STORED)

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
