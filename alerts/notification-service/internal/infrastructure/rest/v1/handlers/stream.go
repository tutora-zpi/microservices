package handlers

import (
	"log"
	"net/http"
)

func (h *SSEHandler) StreamNotifications(w http.ResponseWriter, r *http.Request) {
	h.configureSSEHeaders(w)

	conn, cancel, err := h.prepareSSEConnection(w, r)
	if err != nil {
		h.handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		conn.Cleanup()
		cancel()
	}()

	if err := conn.SendWelcomeMessage(); err != nil {
		log.Printf("Failed to send welcome message to client %s: %v", conn.ClientID, err)
		return
	}

	h.manager.FlushBufferedNotification(conn.ClientID, conn.Channel)

	log.Printf("SSE connection established for client: %s", conn.ClientID)

	conn.HandleEvents()
}
