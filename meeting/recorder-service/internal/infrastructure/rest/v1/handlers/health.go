package handlers

import (
	"net/http"
	"recorder-service/internal/infrastructure/server"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	server.NewResponse(w, "Healthy", http.StatusOK, nil)
}
