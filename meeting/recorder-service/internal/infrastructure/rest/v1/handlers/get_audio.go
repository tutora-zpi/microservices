package handlers

import (
	"log"
	"net/http"
	"recorder-service/internal/domain/dto/request"
	"recorder-service/internal/infrastructure/server"
	"recorder-service/pkg"

	"github.com/gorilla/mux"
)

// GetAudio godoc
// @Summary      Get audio URL
// @Description  Returns a presigned URL for the requested audio file of a meeting
// @Tags         audio
// @Param        meeting_id   path     string  true  "Meeting ID"
// @Param        name         path     string  true  "Audio file name"
// @Success      200  {object}  dto.GetAudioDTO "Presigned URL for audio"
// @Failure      400
// @Router       /api/v1/sessions/audio/{meeting_id}/{name} [get]
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
		return
	}

	server.NewResponse(w, nil, http.StatusOK, *result)
}
