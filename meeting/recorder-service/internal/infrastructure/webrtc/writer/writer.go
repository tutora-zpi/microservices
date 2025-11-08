package writer

import (
	"context"

	"github.com/pion/webrtc/v3"
)

type Writer interface {
	Write(ctx context.Context, userID string, track *webrtc.TrackRemote) error

	GetPath() string
	GetExt() string

	Close()
}

type WriterFactory func(roomID string) (Writer, error)
