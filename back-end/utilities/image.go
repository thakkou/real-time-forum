package utilities

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveImage
func SaveImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Ensure uploads directory exists
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("unable to create uploads directory: %w", err)
	}

	ext := filepath.Ext(fileHeader.Filename) // keep original extension
	newName := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), rand.Intn(10000), ext)
	filePath := filepath.Join("./uploads", newName)

	// Create destination file

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("unable to save file: %w", err)
	}

	// Return the relative URL/path for DB insertion
	return "/uploads/" + newName, nil
}
