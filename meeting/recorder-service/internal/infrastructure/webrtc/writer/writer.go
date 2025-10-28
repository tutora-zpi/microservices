package writer

import (
	"context"

	"github.com/pion/webrtc/v3"
)

type Writer interface {
	Write(ctx context.Context, track *webrtc.TrackRemote)

	GetPath() string
}

type WriterFactory func(roomID, userID string, track *webrtc.TrackRemote) (Writer, error)
