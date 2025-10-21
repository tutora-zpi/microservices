package writer

import (
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

type Writer interface {
	// Saves the data and returns the file path (URL).
	Write(packet *rtp.Packet) error

	// Closes the writer
	Close() error

	GetPath() string
}

type WriterFactory func(track *webrtc.TrackRemote, meetingID string) (Writer, error)
