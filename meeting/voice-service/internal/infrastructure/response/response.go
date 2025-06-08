package response

import (
	"encoding/json"
	"net/http"
)

// Response represents the standard API response format
// @name Response
// @Description Standard API response format containing message, data, and success status
type Response struct {
	// Message contains the response message
	Message string `json:"message"`

	// Data contains the response payload (optional)
	Data any `json:"data,omitempty"`

	// Success indicates whether the operation was successful
	Success bool `json:"success"`
}

// NewResponse creates and sends a standardized JSON response
// Note: This is an internal function and doesn't need Swagger documentation
// as it's not an API endpoint
func NewResponse(w http.ResponseWriter, message string, httpCode int, data any) {
	var success bool = false

	if httpCode < 300 && httpCode >= 200 {
		success = true
	}

	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&Response{
		Message: message,
		Data:    data,
		Success: success,
	}); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
