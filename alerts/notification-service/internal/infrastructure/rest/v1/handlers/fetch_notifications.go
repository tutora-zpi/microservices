package handlers

import (
	"net/http"
	"notification-serivce/internal/domain/query"
	"notification-serivce/internal/infrastructure/server"
	"notification-serivce/pkg"
)

const (
	lastNotificationID string = "last_notification_id"
	limit              string = "limit"
)

func (h *HTTPHandler) FetchNotifications(w http.ResponseWriter, r *http.Request) {
	requesterID, err := ExtractClientID(r)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("Missing user id"), http.StatusBadRequest, nil)
		return
	}

	q := query.NewFetchNotificationsQuery(r.URL.Query().Get(limit), r.URL.Query().Get(lastNotificationID), requesterID)

	res, err := h.bus.HandleQuery(q)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("No notifications yet"), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, res)
}
