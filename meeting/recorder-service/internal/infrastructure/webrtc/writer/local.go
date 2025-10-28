package writer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

type OggWriter struct {
	ow *oggwriter.OggWriter

	path string

	packets []*rtp.Packet

	mutex sync.Mutex
}

// GetPath implements Writer.
func (w *OggWriter) GetPath() string {
	return w.path
}

func (w *OggWriter) flush() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, packet := range w.packets {
		err := w.ow.WriteRTP(packet)
		if err != nil {
			log.Printf("Failed to write RTP packet: %v", err)
		}
	}

	w.packets = nil
}

// Write implements Writer.
func (w *OggWriter) Write(ctx context.Context, track *webrtc.TrackRemote) {
	packetChan := make(chan *rtp.Packet)
	errChan := make(chan error)

	go func() {
		for {
			packet, _, err := track.ReadRTP()
			if err != nil {
				errChan <- err
				return
			}
			packetChan <- packet
		}
	}()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.flush()
			w.ow.Close()
			return
		case <-ticker.C:
			w.flush()
		case packet := <-packetChan:
			w.mutex.Lock()
			w.packets = append(w.packets, packet)
			w.mutex.Unlock()
		case err := <-errChan:
			log.Println("Lost packet:", err)
			return
		}
	}
}

func NewLocalWriter(roomID, userID string, track *webrtc.TrackRemote) (Writer, error) {
	codec := track.Codec()

	basePath, err := dir(roomID)
	if err != nil {
		return nil, err
	}

	savePath := path.Join(basePath, fmt.Sprintf("%s.ogg", userID))

	ow, err := oggwriter.New(savePath, codec.ClockRate, codec.Channels)
	if err != nil {
		return nil, err
	}

	return &OggWriter{ow: ow, path: savePath, mutex: sync.Mutex{}}, nil
}

func dir(roomID string) (string, error) {
	path := path.Join("recordings", roomID)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %s", path)
	}

	return path, nil
}
