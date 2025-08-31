package handlers

import (
	"net/http"
	"notification-serivce/internal/infrastructure/server"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := "Check our docs! Go to /api/v1/docs."
	server.NewResponse(w, &err, http.StatusNotFound, nil)
}
