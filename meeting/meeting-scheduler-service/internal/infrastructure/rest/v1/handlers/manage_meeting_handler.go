package handlers

import (
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/infrastructure/middleware"
	"meeting-scheduler-service/internal/infrastructure/server"
	"meeting-scheduler-service/pkg"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// @title Meeting Scheduler API
	// @version 1.1
	// @description Serivce to requesting meetings in .tutora
	// @host localhost:8003
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
// @Accept json
// @Produce json
// @Param meeting body dto.StartMeetingDTO true "Start Meeting DTO"
// @Success 200 {object} server.Response{data=dto.MeetingDTO} "Meeting details after operation"
// @Failure 400 {object} server.Response "Bad request due to invalid data or DTO type"
// @Router /api/v1/meeting/start [post]
func (m *ManageMeetingHandler) StartMeeting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := middleware.GetDTO(r)
	startDTO, ok := body.(*dto.StartMeetingDTO)
	if !ok {
		server.NewResponse(w, pkg.Ptr("Invalid DTO type"), http.StatusBadRequest, nil)
		return
	}

	meeting, err := m.manager.Start(ctx, *startDTO)
	if err != nil {
		server.NewResponse(w, pkg.Ptr(err.Error()), http.StatusBadRequest, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, meeting)
}

// EndMeeting godoc
// @Summary End a meeting
// @Description Ends a meeting based on the provided DTO
// @Accept json
// @Produce json
// @Param meeting body dto.EndMeetingDTO true "End Meeting DTO"
// @Success 204 {object} server.Response
// @Failure 400 {object} server.Response
// @Router /api/v1/meeting/end [delete]
func (m *ManageMeetingHandler) EndMeeting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := middleware.GetDTO(r)
	endDTO, ok := body.(*dto.EndMeetingDTO)
	if !ok {
		server.NewResponse(w, pkg.Ptr("Invalid body"), http.StatusBadRequest, nil)
		return
	}

	err := m.manager.Stop(ctx, *endDTO)
	if err != nil {
		server.NewResponse(w, pkg.Ptr(err.Error()), http.StatusBadRequest, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusNoContent, nil)
}

// GetActiveMeeting godoc
// @Summary Gets active meeting
// @Description Fetches information about active meeting for given class. "members" will be empty.
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} server.Response{data=dto.MeetingDTO} "Found active meeting"
// @Failure 400 {object} server.Response "Empty class id"
// @Failure 404 {object} server.Response "Not found or not started yet"
// @Router /api/v1/meeting/{class_id} [get]
func (m *ManageMeetingHandler) GetActiveMeeting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	classID := vars["class_id"]
	if classID == "" {
		server.NewResponse(w, pkg.Ptr("Empty class id"), http.StatusBadRequest, nil)
		return
	}

	dto, err := m.manager.ActiveMeeting(ctx, classID)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("Not found or not started yet"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, *dto)
}

// PlanMeeting godoc
// @Summary Plan meeting for the future
// @Description Used to plan meetings and starts meeting automatically at start date
// @Accept json
// @Produce json
// @Success 201 {object} server.Response{data=dto.PlannedMeetingDTO} "Meeting planned successfully!"
// @Failure 400 {object} server.Response "Invalid body | meeting already started | meeting already planned"
// @Router /api/v1/meeting/plan [post]
func (m *ManageMeetingHandler) PlanMeeting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := middleware.GetDTO(r)
	planDTO, ok := body.(*dto.PlanMeetingDTO)
	if !ok {
		server.NewResponse(w, pkg.Ptr("Invalid body"), http.StatusBadRequest, nil)
		return
	}

	meetingDTO, err := m.manager.Plan(ctx, *planDTO)
	if err != nil {
		server.NewResponse(w, pkg.Ptr(err.Error()), http.StatusBadRequest, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusCreated, *meetingDTO)
}

// CancelPlannedMeeting godoc
// @Summary Cancel not started meeting
// @Description Delete from planned meetings, one with provided id
// @Param id path int true "Identifier of planned meeting"
// @Produce json
// @Success 204 {object} server.Response
// @Failure 400 {object} server.Response "Invalid id provided"
// @Router /api/v1/meeting/plan/{id} [delete]
func (m *ManageMeetingHandler) CancelPlannedMeeting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	strID, ok := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strID)
	if !ok || err != nil {
		server.NewResponse(w, pkg.Ptr("Invalid id provided"), http.StatusBadRequest, nil)
		return
	}

	if err := m.manager.CancelPlannedMeeting(ctx, id); err != nil {
		server.NewResponse(w, pkg.Ptr(err.Error()), http.StatusBadRequest, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusNoContent, nil)
}

// GetPlannedMeetings godoc
// @Summary Fetch planned meetings
// @Description Paginated fetch of planned meetings supporting infinite scroll
// @Param class_id path string true "Class ID"
// @Param last_start_timestamp query string false "Last start date timestamp (unix utc format eg. 1760121360)"
// @Param limit query int false "Max number per page, default is 10"
// @Produce json
// @Success 200 {object} server.Response{data=[]dto.PlannedMeetingDTO} "List of planned meetings"
// @Failure 400 {object} server.Response "Invalid parameters"
// @Failure 404 {object} server.Response "No planned meetings found"
// @Router /api/v1/meeting/plan/{class_id} [get]
func (m *ManageMeetingHandler) GetPlannedMeetings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	classID := mux.Vars(r)["class_id"]
	lastStartTimestamp := r.URL.Query().Get("last_start_timestamp")
	limit := r.URL.Query().Get("limit")

	dto, err := dto.NewFetchPlannedMeetingsDTO(classID, limit, lastStartTimestamp)
	log.Print(*dto)
	if err != nil {
		server.NewResponse(w, pkg.Ptr(err.Error()), http.StatusBadRequest, nil)
		return
	}

	plannedMeetings, err := m.manager.GetPlannedMeetings(ctx, *dto)

	if err != nil {
		server.NewResponse(w, pkg.Ptr(err.Error()), http.StatusBadRequest, nil)
		return
	}

	if len(plannedMeetings) < 1 {
		server.NewResponse(w, pkg.Ptr("Not found matched planned meetings"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, plannedMeetings)
}
