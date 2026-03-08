package helper

import (
	"forum/database"
)

func GetUserIDFromCookie(cookie string) (int, error) {

	var userID int
	err := database.Database.QueryRow(`
		SELECT user_id
		FROM SESSIONS
		WHERE id = ? AND expires_at > DATETIME('now')
	`, cookie).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}