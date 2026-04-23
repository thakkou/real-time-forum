package utilities

import (
	"bytes"
	"html/template"
	"net/http"
)

// RenderTemplate
func RenderTemplate(w http.ResponseWriter, status int, tmpl string, data any) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Template error")
		return
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

// HandleError renders an error page with the given status and message
func HandleError(w http.ResponseWriter, code int, message string) {
	// 	w.WriteHeader(code)

	t, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Critical template error", http.StatusInternalServerError) // "Internal Server Error"
		return
	}

	data := struct {
		Code    int
		Message string
	}{
		Code:    code,
		Message: message,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		// log.Printf("error template execute error: %v", err)
		return
	}
	w.WriteHeader(code)
	buf.WriteTo(w)
}
