package handlers

import (
	"net/http"
	"recorder-service/internal/app/interfaces/service"
	// @title Recording Serivce API
	// @version 1.0
	// @description Service responsible for recording meetings
	// @host localhost:8050
)

type Handler interface {
	FetchSessions(w http.ResponseWriter, r *http.Request)
	GetAudio(w http.ResponseWriter, r *http.Request)
	NotFound(w http.ResponseWriter, r *http.Request)
	NotFoundHandler(w http.ResponseWriter, r *http.Request)
	IsAuth(next http.Handler) http.Handler
}

type handlerImpl struct {
	voiceService service.VoiceSessionService
}

func (h *handlerImpl) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/v1/docs", http.StatusSeeOther)
}

func NewHandler(voiceService service.VoiceSessionService) Handler {
	return &handlerImpl{voiceService: voiceService}
}
