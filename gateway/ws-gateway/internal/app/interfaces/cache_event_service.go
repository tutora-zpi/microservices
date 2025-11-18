package interfaces

import (
	"context"
	wsevent "ws-gateway/internal/domain/ws_event"
)

type CacheEventService interface {
	PushRecentEvent(ctx context.Context, event wsevent.SocketEventWrapper, roomID string) error
	GetLastEventsData(ctx context.Context, roomID string) ([][]byte, error)
	MakeSnapshot(ctx context.Context, event wsevent.SocketEventWrapper, roomID string) error
	GetSnapshot(ctx context.Context, keroomIDy string) ([]byte, error)
}
