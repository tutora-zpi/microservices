package handlers

import (
	"net/http"
	"voice-service/internal/infrastructure/response"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	response.NewResponse(w, "Healthy", http.StatusOK, nil)
}
