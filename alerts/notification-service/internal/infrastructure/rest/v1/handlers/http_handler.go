package handlers

import (
	"encoding/json"
	"net/http"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/infrastructure/server"
	"notification-serivce/pkg"
	// @title Notifications Serivce API
	// @version 1.0
	// @description Service responsible for notification delivery
	// @host localhost:8888
)

type HTTPHandler struct {
	serivce interfaces.NotificationSerivce
}

const (
	lastNotificationID string = "last_notification_id"
	limit              string = "limit"
)

func NewHTTPHandler(serivce interfaces.NotificationSerivce) *HTTPHandler {
	return &HTTPHandler{serivce: serivce}
}

// FetchNotifications godoc
// @Summary List with user's notifications
// @Description Paginated reponse adjusted for infinite scroll
// @Tags notifications
// @Produce json
// @Param limit query int false "Number of notifications to display, limit range (0, 20] if number is invalid then limit is set to 10"
// @Param last_notification_id query string false "ID of the last notification for pagination, not required during first fetch."
// @Success 200 {object} server.Response{data=dto.NotificationDTO} "Notification's data"
// @Failure 400 {object} server.Response "Missing user id"
// @Failure 404 {object} server.Response "No notifications yet"
// @Router /api/v1/notification [get]
func (h *HTTPHandler) FetchNotifications(w http.ResponseWriter, r *http.Request) {
	requesterID, err := ExtractClientID(r)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("Missing user id"), http.StatusBadRequest, nil)
		return
	}

	req := dto.NewFetchNotificationsDTO(r.URL.Query().Get(limit), r.URL.Query().Get(lastNotificationID), requesterID)

	res, err := h.serivce.FetchNotifications(req)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("No notifications yet"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, res)
}

// DeleteNotifications godoc
// @Summary Deletes notifications
// @Description Removes one or more notifications for a given user
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body dto.DeleteNotificationsDTO true "IDs of notifications to delete"
// @Success 204 "No Content"
// @Failure 400 {object} server.Response "Invalid body or Missing user id"
// @Failure 404 {object} server.Response "Notification not found"
// @Router /api/v1/notification [delete]
func (h *HTTPHandler) DeleteNotifications(w http.ResponseWriter, r *http.Request) {
	clientID, err := ExtractClientID(r)
	if err != nil {
		// server.NewResponse(w, pkg.Ptr("Missing user id"), http.StatusBadRequest, nil)
		// return
		clientID = "91e56313-e147-428a-844e-734f0e869f6f"
	}

	defer r.Body.Close()

	var req dto.DeleteNotificationsDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.NewResponse(w, pkg.Ptr("Invalid body"), http.StatusBadRequest, nil)
		return
	}

	if err := h.serivce.DeleteNotifications(&req, clientID); err != nil {
		server.NewResponse(w, pkg.Ptr("Notification not found"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusNoContent, nil)
}
