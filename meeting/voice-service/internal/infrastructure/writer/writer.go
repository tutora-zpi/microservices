package writer

import (
	"github.com/pion/rtp"
)

type Writer interface {
	// Saves the data and returns the file path (URL).
	Write(packet *rtp.Packet) error

	// Closes the writer
	Close() error

	GetPath() string
}
