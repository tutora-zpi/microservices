package general

import (
	"context"
	"fmt"
	"signaling-service/internal/app/interfaces"
	"time"

	"github.com/gorilla/websocket"
)

type heartbeatHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (u *heartbeatHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()
			nowStr := fmt.Sprintf("%d", now)
			u.hubManager.EmitGlobal(websocket.PingMessage, []byte(nowStr))
		}
	}

	return nil
}

func NewHeartbeatHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &userLeftHandler{hubManager: hubManager}
}
