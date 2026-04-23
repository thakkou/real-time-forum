package handlers

import (
	"log"
	"net/http"

	"forum/models"
	"forum/utilities"
)

// Logout
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		utilities.HandleError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodPost {
		utilities.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil { // http.ErrNoCookie
		return
	}

	err = models.DeleteSession(cookie.Value)
	// + need to remove cookie from storage
	if err != nil {
		log.Println(err)
	}

	http.SetCookie(w, &http.Cookie{ // delete cookie ------------------- TODO: function already exists
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther) // or to login
}
