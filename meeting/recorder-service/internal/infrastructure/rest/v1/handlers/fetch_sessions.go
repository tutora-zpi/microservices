package handlers

import (
	"net/http"
	"recorder-service/internal/domain/dto/request"
	"recorder-service/internal/infrastructure/server"
	"recorder-service/pkg"

	"github.com/gorilla/mux"
)

// FetchSessions godoc
// @Summary      Fetch sessions metadata
// @Description  Returns a list of voice session metadata for a given meeting. Supports pagination with lastFetchedId and limit.
// @Tags         sessions
// @Param        meeting_id    path     string  true   "Meeting ID"
// @Param        lastFetchedId query    string  false  "ID of the last fetched session for pagination"
// @Param        limit         query    int     false  "Maximum number of sessions to return"
// @Success      200  {array}   dto.VoiceSessionMetadataDTO "List of session metadata"
// @Failure      400
// @Failure      404
// @Router       /api/v1/sessions/{meeting_id} [get]
func (h *handlerImpl) FetchSessions(w http.ResponseWriter, r *http.Request) {
	meetingID, ok := mux.Vars(r)["meeting_id"]
	if !ok {
		server.NewResponse(w, pkg.Ptr("No class id"), http.StatusBadRequest, nil)
		return
	}

	ctx := r.Context()

	lastFetchedID := r.URL.Query().Get("lastFetchedId")
	limit := r.URL.Query().Get("limit")

	req := request.NewFetchSessions(meetingID, lastFetchedID, limit)

	result, err := h.voiceService.GetSessions(ctx, *req)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("Not found matching sessions metadata"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, result)
}
