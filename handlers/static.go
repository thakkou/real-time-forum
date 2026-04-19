package handlers

import (
	"net/http"
	"os"
)

// Static serves static files (css, js, images...)
func Static(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/assets" ||
		r.URL.Path == "/assets/" ||
		r.URL.Path == "/uploads" ||
		r.URL.Path == "/uploads/" {
		HandleError(w, http.StatusNotFound, "Resource Not Found")
		return
	}

	var filePath string = r.URL.Path[1:]
	if _, err := os.Stat(filePath); err != nil {
		HandleError(w, http.StatusNotFound, "Resource Not Found")
		return
	}
	http.ServeFile(w, r, filePath)
}
