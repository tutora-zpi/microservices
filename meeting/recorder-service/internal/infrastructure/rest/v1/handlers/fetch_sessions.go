package handlers

import (
	"net/http"
	"recorder-service/internal/domain/dto/request"
	"recorder-service/internal/infrastructure/server"
	"recorder-service/pkg"

	"github.com/gorilla/mux"
)

// FetchSessions implements Handler.
func (h *handlerImpl) FetchSessions(w http.ResponseWriter, r *http.Request) {
	classID, ok := mux.Vars(r)["class_id"]
	if !ok {
		server.NewResponse(w, pkg.Ptr("No class id"), http.StatusBadRequest, nil)
		return
	}

	ctx := r.Context()

	lastFetchedMeetingID := r.URL.Query().Get("last_fetched_meeting_id")
	limit := r.URL.Query().Get("limit")

	req := request.NewFetchSessions(classID, lastFetchedMeetingID, limit)

	result, err := h.voiceService.GetSessions(ctx, *req)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("Not found matching sessions metadata"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, result)
}
