package handlers

import (
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/infrastructure/middleware"
	"meeting-scheduler-service/internal/infrastructure/response"
	"net/http"
	// @title Meeting Scheduler API
	// @version 1.0
	// @description Serivce to requesting meetings in .tutora
	// @host default localhost:8080
	// @BasePath /api/v1
)

type ManageMeetingHandler struct {
	manager interfaces.ManageMeeting
}

func NewManageMeetingHandler(m interfaces.ManageMeeting) ManageMeetingHandler {
	return ManageMeetingHandler{
		manager: m,
	}
}

// StartMeeting godoc
// @Summary Start a meeting
// @Description Starts a meeting based on the provided DTO
// @Tags meetings
// @Accept json
// @Produce json
// @Param meeting body dto.StartMeetingDTO true "Start Meeting DTO"
// @Success 200 {object} response.Response "Meeting details after operation"
// @Failure 400 {object} response.Response "Bad request due to invalid data or DTO type"
// @Failure 405 {object} response.Response "Method not allowed (only POST supported)"
// @Router /api/v1/meeting/start [post]
func (m *ManageMeetingHandler) StartMeeting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.NewResponse(w, "Use Post method", http.StatusMethodNotAllowed, nil)
		return
	}

	body := middleware.GetDTO(r)
	startDTO, ok := body.(*dto.StartMeetingDTO)
	if !ok {
		response.NewResponse(w, "invalid DTO type", http.StatusBadRequest, nil)
		return
	}

	log.Println("Requested starting meeting")
	meeting, err := m.manager.Start(*startDTO)
	if err != nil {
		response.NewResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	response.NewResponse(w, "Meeting started successfully!", http.StatusOK, meeting)
}

// EndMeeting godoc
// @Summary End a meeting
// @Description Ends a meeting based on the provided DTO
// @Tags meetings
// @Accept json
// @Produce json
// @Param meeting body dto.EndMeetingDTO true "End Meeting DTO"
// @Success 200 {object} response.Response "Meeting details after operation"
// @Failure 400 {object} response.Response "Bad request due to invalid data or DTO type"
// @Failure 405 {object} response.Response "Method not allowed (only POST supported)"
// @Router /api/v1/meeting/end [post]
func (m *ManageMeetingHandler) EndMeeting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.NewResponse(w, "Use Post method", http.StatusMethodNotAllowed, nil)
		return
	}

	body := middleware.GetDTO(r)
	endDTO, ok := body.(*dto.EndMeetingDTO)
	if !ok {
		response.NewResponse(w, "invalid DTO type", http.StatusBadRequest, nil)
		return
	}

	log.Println("Requested stopping meeting")
	meeting, err := m.manager.Stop(*endDTO)
	if err != nil {
		response.NewResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	response.NewResponse(w, "Meeting ended successfully!", http.StatusOK, meeting)
}
