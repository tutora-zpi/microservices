package writer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

type OggWriter struct {
	basePath string
	done     chan struct{}
	ext      string
}

// Close implements Writer.
func (w *OggWriter) Close() {
	select {
	case <-w.done:
		return
	default:
		close(w.done)
	}
}

// GetPath implements Writer.
func (w *OggWriter) GetPath() string {
	return w.basePath
}

// GetExt implements Writer.
func (w *OggWriter) GetExt() string {
	return w.ext
}

// Write implements Writer.
func (w *OggWriter) Write(ctx context.Context, userID string, track *webrtc.TrackRemote) error {
	filename := path.Join(w.basePath, fmt.Sprintf("%s.%s", userID, w.ext))

	ow, err := oggwriter.New(filename, 48000, track.Codec().Channels)
	if err != nil {
		log.Printf("Failed to create writer: %v", err)
		return fmt.Errorf("failed to create writer")
	}

	ticker := time.NewTicker(time.Second)

	go func() {
		defer ow.Close()

		for {
			packet, _, err := track.ReadRTP()
			if err != nil {
				log.Printf("Track from %s ended: %v", userID, err)
				return
			}

			if err := ow.WriteRTP(packet); err != nil {
				log.Printf("Failed to write RTP for %s: %v", userID, err)
				return
			}

			select {
			case <-ticker.C:
				log.Printf("Receiving packets from: %s", userID)
			case <-ctx.Done():
			case <-w.done:
				return
			default:
			}
		}
	}()

	return nil
}

func NewLocalWriter(roomID string) (Writer, error) {

	path, err := setupBasePath(roomID)
	if err != nil {
		return nil, err
	}

	return &OggWriter{
		basePath: path,
		done:     make(chan struct{}),
		ext:      "ogg",
	}, nil
}

func setupBasePath(roomID string) (string, error) {
	p := path.Join("recordings", roomID)
	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %s", p)
	}

	return p, nil
}
