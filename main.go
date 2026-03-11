package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/routing"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	http.HandleFunc("/static/", handlers.HandleStatic)
	routing.RegisterRoutes()

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
