package handler

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func writeValidationError(w http.ResponseWriter, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  "validation failed", 
		"fields": fields,
	}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}