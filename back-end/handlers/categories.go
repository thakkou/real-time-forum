package handlers

import (
	"fmt"
	"forum/database"
)

func GetCategoriesByPost(postId int) ([]string, error) {
	var categories []string

	rows, err := database.Database.Query(`
		SELECT c.name
		FROM category c
		JOIN post_category pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
		ORDER BY c.name
	`, postId)
	if err != nil {
		return nil, fmt.Errorf("GetCategoriesByPost error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("GetCategoriesByPost scan error: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCategoriesByPost rows error: %v", err)
	}

	return categories, nil
}
