package util

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondJSON(statusCode int, responseBody interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		log.Printf("Error writing the response: %v", err)
	}
}
