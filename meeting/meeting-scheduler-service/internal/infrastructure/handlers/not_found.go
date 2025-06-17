package handlers

import (
	"meeting-scheduler-service/internal/infrastructure/response"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.NewResponse(w, "Check our docs! Go to /api/v1/docs.", http.StatusNotFound, nil)
}
