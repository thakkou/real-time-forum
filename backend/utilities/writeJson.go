package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Println("status code", statusCode, "message", message)
	json.NewEncoder(w).Encode(Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}
