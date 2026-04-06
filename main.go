package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"forum/database"
	"forum/handlers"
	"forum/routes"
)

func main() {
	if err := InitEnviron(); err != nil {
		log.Fatalf("Environ initialization failed: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	http.HandleFunc("/static/", handlers.Static)
	http.HandleFunc("/", handlers.Forum)

	routes.RegisterRoutes()

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func InitEnviron() error {
	bytes, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	handlers.GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
	handlers.GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
	handlers.GITHUB_CLIENT_ID = os.Getenv("GITHUB_CLIENT_ID")
	handlers.GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
	return nil
}
