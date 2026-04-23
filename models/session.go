package models

import "forum/database"

func DeleteSession(sessionId string) error {
	// return database.Database.QueryRow(query).Err()
	// db.exec vs db.queryrow in golang sqlite
	// queryrow not working with delete statement
	_, err := database.Database.Exec(
		"DELETE FROM sessions WHERE id = ?",
		sessionId) // returns result
	return err
}

func GetUserIDFromCookie(sessionID string) (int, error) {
	var userID int
	err := database.Database.QueryRow(`
		SELECT user_id
		FROM SESSIONS
		WHERE id = ?
	`, sessionID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
