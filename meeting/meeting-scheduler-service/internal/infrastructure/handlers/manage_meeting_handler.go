package handlers

import (
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/infrastructure/middleware"
	"meeting-scheduler-service/internal/infrastructure/response"
	"net/http"

	"github.com/gorilla/mux"
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
// @Tags start meetings
// @Accept json
// @Produce json
// @Param meeting body dto.StartMeetingDTO true "Start Meeting DTO"
// @Success 200 {object} response.Response{data=dto.MeetingDTO} "Meeting details after operation"
// @Failure 400 {object} response.Response "Bad request due to invalid data or DTO type"
// @Failure 405 {object} response.Response "Method not allowed (only POST supported)"
// @Router /api/v1/meeting/start [post]
// @Router /api/v1/meeting/start [put]
func (m *ManageMeetingHandler) StartMeeting(w http.ResponseWriter, r *http.Request) {
	body := middleware.GetDTO(r)
	startDTO, ok := body.(*dto.StartMeetingDTO)
	if !ok {
		response.NewResponse(w, "Invalid DTO type", http.StatusBadRequest, nil)
		return
	}

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
// @Tags ending meetings
// @Accept json
// @Produce json
// @Param meeting body dto.EndMeetingDTO true "End Meeting DTO"
// @Success 200 {object} response.Response{data=dto.MeetingDTO}
// @Failure 400 {object} response.Response
// @Router /api/v1/meeting/end [post]
// @Router /api/v1/meeting/end [delete]
func (m *ManageMeetingHandler) EndMeeting(w http.ResponseWriter, r *http.Request) {
	body := middleware.GetDTO(r)
	endDTO, ok := body.(*dto.EndMeetingDTO)
	if !ok {
		response.NewResponse(w, "Invalid body", http.StatusBadRequest, nil)
		return
	}

	err := m.manager.Stop(*endDTO)
	if err != nil {
		response.NewResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	response.NewResponse(w, "Meeting ended successfully!", http.StatusOK, nil)
}

// GetActiveMeeting godoc
// @Summary Gets active meeting
// @Description Fetches information about active meeting for given class. "members" will be empty.
// @Tags meetings class
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} response.Response{data=dto.MeetingDTO} "Found active meeting"
// @Failure 400 {object} response.Response "Empty class id"
// @Failure 404 {object} response.Response "Not found or not started yet"
// @Failure 405 {object} response.Response "Method not allowed (only GET supported)"
// @Router /api/v1/meeting/{class_id} [get]
func (m *ManageMeetingHandler) GetActiveMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classID := vars["class_id"]
	if classID == "" {
		response.NewResponse(w, "Empty class id", http.StatusBadRequest, nil)
		return
	}

	dto, err := m.manager.ActiveMeeting(classID)
	if err != nil {
		response.NewResponse(w, "Not found or not started yet", http.StatusNotFound, nil)
		return
	}

	response.NewResponse(w, "Found active meeting", http.StatusOK, *dto)
}
