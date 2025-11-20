package handlers

import (
	"log"
	"net/http"
	"recorder-service/internal/domain/dto/request"
	"recorder-service/internal/infrastructure/server"
	"recorder-service/pkg"

	"github.com/gorilla/mux"
)

// GetAudio implements Handler.
func (h *handlerImpl) GetAudio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meetingID, ok := mux.Vars(r)["meeting_id"]
	if !ok {
		server.NewResponse(w, pkg.Ptr("No meetings id"), http.StatusBadRequest, nil)
		return
	}

	audioName, ok := mux.Vars(r)["name"]
	if !ok {
		server.NewResponse(w, pkg.Ptr("No meetings id"), http.StatusBadRequest, nil)
		return
	}

	req := request.GetAudioRequest{AudioName: audioName, MeetingID: meetingID}

	result, err := h.voiceService.GetAudio(ctx, req)
	if err != nil {
		log.Printf("Something went wrong during fetching audio url: %v", err)
		server.NewResponse(w, pkg.Ptr(err.Error()), http.StatusBadRequest, nil)
	}

	server.NewResponse(w, nil, http.StatusOK, *result)
}
