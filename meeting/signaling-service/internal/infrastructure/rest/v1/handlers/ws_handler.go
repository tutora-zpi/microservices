package handlers

import (
	"fmt"
	"net/http"
	"signaling-service/internal/infrastructure/ws"
)

// WebSocketHandler implements Handlable.
func (h *handlers) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	ctx := r.Context()
	// id := r.Context().Value(id).(string)
	id := "lukasz"

	client := ws.NewClient(id, conn)

	h.hub.AddGlobalMember(client)

	// go client.Listen(ctx, h.WithAuth(h.dispatcher.HandleEvent))
	go client.Listen(ctx, h.dispatcher.HandleEvent)
}
