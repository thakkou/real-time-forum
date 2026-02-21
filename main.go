package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/handlers"
)

func main() {
	database.Init()

	http.HandleFunc("/", handlers.Forum)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	// http.HandleFunc("/logout", )

	http.HandleFunc("/static/styles.css", handlers.CssHandler)

	fmt.Println("Server running on http://0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
