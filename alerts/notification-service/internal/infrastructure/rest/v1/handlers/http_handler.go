package handlers

import (
	"encoding/json"
	"net/http"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/requests"
	"notification-serivce/internal/infrastructure/server"
	"notification-serivce/pkg"
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

func (h *HTTPHandler) FetchNotifications(w http.ResponseWriter, r *http.Request) {
	requesterID, err := ExtractClientID(r)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("Missing user id"), http.StatusBadRequest, nil)
		return
	}

	req := requests.NewFetchNotificationsRequest(r.URL.Query().Get(limit), r.URL.Query().Get(lastNotificationID), requesterID)

	res, err := h.serivce.FetchNotifications(req)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("No notifications yet"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, res)
}

func (h *HTTPHandler) DeleteNotifications(w http.ResponseWriter, r *http.Request) {
	clientID, err := ExtractClientID(r)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("Missing user id"), http.StatusBadRequest, nil)
		return
	}

	defer r.Body.Close()

	var req requests.DeleteNotificationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.NewResponse(w, pkg.Ptr("Invalid body"), http.StatusBadRequest, nil)
		return
	}

	if err := h.serivce.DeleteNotifications(&req, clientID); err != nil {
		server.NewResponse(w, pkg.Ptr("No notifications yet"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusNoContent, nil)
}
