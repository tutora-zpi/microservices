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
	router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	api := router.PathPrefix("/api/v1/meeting").Subrouter()
	api.Handle("/start", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.StartMeeting))))
	api.Handle("/end", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.EndMeeting))))

	return router
}
