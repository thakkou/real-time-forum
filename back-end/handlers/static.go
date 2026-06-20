package handlers

import (
	"net/http"
	"os"

	"forum/utilities"
)

// Static serves static files (css, js, images..)
func Static(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/assets" ||
		r.URL.Path == "/assets/" ||
		r.URL.Path == "/uploads" ||
		r.URL.Path == "/uploads/" {
		utilities.WriteJSON(w, http.StatusNotFound, "Resource Not Found", nil)
		return
	}

	var filePath string = r.URL.Path[1:]
	if _, err := os.Stat(filePath); err != nil {
		utilities.WriteJSON(w, http.StatusNotFound, "Resource Not Found", nil)
		return
	}
	http.ServeFile(w, r, filePath)
}
