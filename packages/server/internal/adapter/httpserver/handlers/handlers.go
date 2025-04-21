package handlers

import (
	"net/http"

	"log"
)

func RespondWithError(w http.ResponseWriter, r *http.Request, status int, message string, err error) {
	log.Printf("Error: %v, Status: %d, Message: %s", err, status, message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error": "` + message + `"}`))
}
