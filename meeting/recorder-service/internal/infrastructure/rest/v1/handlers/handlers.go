package handlers

import (
	"net/http"
	"recorder-service/internal/app/interfaces/service"
)

type Handler interface {
	FetchSessions(w http.ResponseWriter, r *http.Request)
	GetAudio(w http.ResponseWriter, r *http.Request)
	NotFound(w http.ResponseWriter, r *http.Request)
	IsAuth(next http.Handler) http.Handler
}

type handlerImpl struct {
	voiceService service.VoiceSessionService
}

func NewHandler(voiceService service.VoiceSessionService) Handler {
	return &handlerImpl{voiceService: voiceService}
}
