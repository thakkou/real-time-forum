package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/routes"
	"forum/utilities"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	utilities.WriteJSON(w, 200, "server is healty", nil)
	fmt.Println("healt")
	return
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Allow common methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// if err := InitEnviron(); err != nil {
	// 	log.Fatalf("Environ initialization failed: %v", err)
	// }

	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	// this health check for server

	http.HandleFunc("/health", healthHandler)

	// this websokets
	http.HandleFunc("/ws", handlers.HandlerWs)
	// test ws
	http.HandleFunc("/ws/test", handlers.TestBroadcast)
	http.HandleFunc("/assets/", handlers.Static)
	http.HandleFunc("/uploads/", handlers.Static)

	// http.HandleFunc("/", handlers.Forum)

	routes.RegisterRoutes()

	log.Println("Server running on http://localhost:8080")
	handler := corsMiddleware(http.DefaultServeMux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

// func InitEnviron() error {
// 	bytes, err := os.ReadFile(".env")
// 	if err != nil {
// 		return err
// 	}
// 	lines := strings.Split(string(bytes), "\n")
// 	for _, line := range lines {
// 		if len(line) == 0 || strings.HasPrefix(line, "#") {
// 			continue
// 		}

// 		// Split key and value
// 		parts := strings.SplitN(line, "=", 2)
// 		if len(parts) != 2 {
// 			continue
// 		}

// 		key := strings.TrimSpace(parts[0])
// 		value := strings.TrimSpace(parts[1])
// 		os.Setenv(key, value)
// 	}

// 	handlers.GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
// 	handlers.GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
// 	handlers.GITHUB_CLIENT_ID = os.Getenv("GITHUB_CLIENT_ID")
// 	handlers.GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
// 	return nil
// }
