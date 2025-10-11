package rest

import (
	"meeting-scheduler-service/internal/infrastructure/middleware"
	"meeting-scheduler-service/internal/infrastructure/rest/v1/handlers"
	"net/http"

	_ "meeting-scheduler-service/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(meetingHandler handlers.ManageMeetingHandler) *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	api.Handle("/docs", http.RedirectHandler("/api/v1/docs/", http.StatusSeeOther))

	meeting := api.PathPrefix("/meeting").Subrouter()

	meeting.Handle("/start", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.StartMeeting)))).Methods(http.MethodPost, http.MethodPut)
	meeting.Handle("/end", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.EndMeeting)))).Methods(http.MethodPost, http.MethodDelete)
	meeting.Handle("/{class_id}", middleware.IsAuth(http.HandlerFunc(meetingHandler.GetActiveMeeting))).Methods(http.MethodGet)

	plan := meeting.PathPrefix("/plan").Subrouter()

	plan.Handle("", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.PlanMeeting)))).Methods(http.MethodPost)
	plan.Handle("/{id}/cancel", middleware.IsAuth(http.HandlerFunc(meetingHandler.CancelPlannedMeeting))).Methods(http.MethodDelete)
	plan.Handle("/{class_id}", middleware.IsAuth(http.HandlerFunc(meetingHandler.GetPlannedMeetings))).Methods(http.MethodGet)

	return router
}
