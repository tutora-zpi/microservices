package rest

import (
	"meeting-scheduler-service/internal/infrastructure/handlers"
	"meeting-scheduler-service/internal/infrastructure/middleware"
	"net/http"

	_ "meeting-scheduler-service/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(meetingHandler handlers.ManageMeetingHandler) *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	router.PathPrefix("/api/v1/docs/").Handler(httpSwagger.WrapHandler)

	api := router.PathPrefix("/api/v1").Subrouter()
	// api.Handle("/start", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.StartMeeting))))
	// api.Handle("/end", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.EndMeeting))))

	api.Handle("/meeting/start", middleware.Validate(http.HandlerFunc(meetingHandler.StartMeeting))).Methods(http.MethodPost, http.MethodPut)
	api.Handle("/meeting/end", middleware.Validate(http.HandlerFunc(meetingHandler.EndMeeting))).Methods(http.MethodPost, http.MethodDelete)
	api.Handle("/meeting/{class_id}", http.HandlerFunc(meetingHandler.GetActiveMeeting)).Methods(http.MethodGet)

	return router
}
