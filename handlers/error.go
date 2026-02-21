package handlers

import (
	"bytes"
	"html/template"
	"net/http"
)

// HandleError renders an error page with the given status and message
func HandleError(w http.ResponseWriter, status int, message string) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
	data := struct {
		Message string
		Status  int
	}{
		Message: message,
		Status:  status,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}
