package interfaces

import (
	"context"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/recorder"
)

type CacheEventService interface {
	PushRecentEvent(ctx context.Context, event wsevent.SocketEventWrapper, roomID string) error
	GetLastEventsData(ctx context.Context, roomID string) ([][]byte, error)
	MakeSnapshot(ctx context.Context, event wsevent.SocketEventWrapper, roomID string) error
	GetSnapshot(ctx context.Context, keyRoomID string) ([]byte, error)
	IsMeetingRecorded(context.Context, string) (*recorder.RecordRequestedWSEvent, error)
	SetMeetingIsRecorded(ctx context.Context, keyRoomID string, evt recorder.RecordRequestedWSEvent) error
	RemoveMeetingFromPool(ctx context.Context, keyRoomID string) error
}
