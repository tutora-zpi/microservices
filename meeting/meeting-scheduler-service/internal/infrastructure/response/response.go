package response

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Success bool   `json:"sucess"`
}

func NewResponse(w http.ResponseWriter, message string, httpCode int, data any) {
	var success bool = false

	if httpCode < 300 && httpCode >= 200 {
		success = true
	}

	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&response{
		Message: message,
		Data:    data,
		Success: success,
	}); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
