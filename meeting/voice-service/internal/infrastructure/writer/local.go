package writer

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

type LocalWriter struct {
	ow *oggwriter.OggWriter

	path string
}

// Close implements Writer.
func (w *LocalWriter) Close() error {
	return w.ow.Close()
}

// WriteLocally implements Writer.
func (w *LocalWriter) Write(packet *rtp.Packet) error {
	return w.ow.WriteRTP(packet)
}

func NewLocalWriter(track *webrtc.TrackRemote, meetingID string) (Writer, error) {
	codec := track.Codec()

	path := generateNewPath(meetingID)

	ow, err := oggwriter.New(path, codec.ClockRate, codec.Channels)
	if err != nil {
		return nil, err
	}

	return &LocalWriter{ow: ow, path: path}, nil
}

func (w *LocalWriter) GetPath() string {
	return w.path
}

// Generates new path for voice recording files if the structure does not exists creates it.
// Please provide meeting id as a filename
func generateNewPath(filename string) string {
	src := path.Join("voice", "recordings")

	if !pathExists(src) {
		creatDirs(src)
	}

	name := fmt.Sprintf("%s_%d.ogg", filename, time.Now().UnixNano()/int64(time.Millisecond))
	return path.Join(src, name)
}

func pathExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func creatDirs(src string) {
	if err := os.MkdirAll(src, 0755); err != nil {
		log.Printf("Failed to create directory %s: %v", src, err)
	}
}
