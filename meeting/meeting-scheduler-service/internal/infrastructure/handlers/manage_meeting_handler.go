package handlers

import (
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/infrastructure/middleware"
	"meeting-scheduler-service/internal/infrastructure/response"
	"net/http"
)

type ManageMeetingHandler struct {
	manager interfaces.ManageMeeting
}

func NewManageMeetingHandler(m interfaces.ManageMeeting) ManageMeetingHandler {
	return ManageMeetingHandler{
		manager: m,
	}
}

func (m *ManageMeetingHandler) Handler(w http.ResponseWriter, r *http.Request) {
	var err error
	var meeting *dto.MeetingDTO

	body := middleware.GetDTO(r)
	switch v := body.(type) {
	case *dto.StartMeetingDTO:
		log.Println("Requested starting meeting")
		meeting, err = m.manager.Start(*v)
	case *dto.EndMeetingDTO:
		log.Println("Requested stopping meeting")
		meeting, err = m.manager.Stop(*v)
	default:
		response.NewResponse(w, "unexpected dto type", http.StatusBadRequest, nil)
		return
	}

	if err != nil {
		response.NewResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	response.NewResponse(w, "Opertaion has been excetured successfully!", http.StatusOK, meeting)
}
