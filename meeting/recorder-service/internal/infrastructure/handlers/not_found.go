package handlers

import (
	"net/http"
	"recorder-service/internal/infrastructure/response"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.NewResponse(w, "Check out docs! Go to /api/v1/docs", http.StatusNotFound, nil)
}
