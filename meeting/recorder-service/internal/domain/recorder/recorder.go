package recorder

import (
	"context"
	"recorder-service/internal/domain/client"
	"recorder-service/internal/infrastructure/webrtc/writer"
)

type Recorder interface {
	StartRecording(ctx context.Context, roomID string) error
	StopRecording(ctx context.Context, roomID string) ([]string, error)
	AddRecordedClient(ctx context.Context, writerFactory writer.WriterFactory, roomID, userID string, client client.Client) error
	RegisterNewRoom(ctx context.Context, roomID string) error
}
