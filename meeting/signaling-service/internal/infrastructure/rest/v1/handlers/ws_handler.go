package handlers

import (
	"context"
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

	/// TOODO CHANGE it
	id := r.URL.Query().Get("id")

	client := ws.NewClient(id, conn)

	h.hub.AddGlobalMember(client)

	backgroundCtx := createBackgroundCtx(ctx)

	// go client.Listen(backgroundCtx, h.WithAuth(h.dispatcher.HandleEvent))

	go client.Listen(backgroundCtx, h.dispatcher.HandleEvent)
}

func createBackgroundCtx(requestCtx context.Context) context.Context {
	userID, _ := requestCtx.Value(ID).(string)
	token, _ := requestCtx.Value(Token).(string)

	background := context.Background()
	background = context.WithValue(background, ID, userID)
	background = context.WithValue(background, Token, token)

	return background
}
