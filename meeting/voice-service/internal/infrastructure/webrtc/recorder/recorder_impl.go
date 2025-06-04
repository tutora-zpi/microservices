package recorder

import (
	"log"
	"sync"
	"voice-service/internal/app/interfaces"
	"voice-service/internal/infrastructure/writer"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

type recorderImpl struct {
	conn   *webrtc.PeerConnection
	writer writer.Writer

	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewRecorder(writer writer.Writer, cfg *webrtc.Configuration) interfaces.Recorder {
	if cfg == nil {
		cfg = &webrtc.Configuration{}
	}

	peerConnection, err := webrtc.NewPeerConnection(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &recorderImpl{
		conn: peerConnection,
	}
}

// Starts recording audio tracks and saves them to the specified file.
func (r *recorderImpl) StartRecording() {
	r.stopChan = make(chan struct{})

	r.conn.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		r.handleBuffer(track)
	})
}

// Stops recording and waits for all goroutines to finish.
func (r *recorderImpl) StopRecording() {
	if r.stopChan != nil {
		close(r.stopChan)
		r.wg.Wait()
		log.Println("Recording finished!")
		return
	}

	log.Println("Recording already stopped or not started (stopChan is nil).")
}

func (r *recorderImpl) handleBuffer(track *webrtc.TrackRemote) {
	if track.Kind() != webrtc.RTPCodecTypeAudio {
		log.Println("Wrong type:", track.Kind())
		return
	}

	log.Println("Reading audio:", track.ID(), "type:", track.Kind())
	r.wg.Add(1)
	go r.receiveAudio(track)
}

func (r *recorderImpl) receiveAudio(track *webrtc.TrackRemote) {
	defer r.wg.Done()

	defer func() {
		if err := r.writer.Close(); err != nil {
			log.Println("Failed to close ogg writer:", err)
		}
	}()

	buf := make([]byte, 1400)
	for {
		select {
		case <-r.stopChan:
			log.Println("Receving audio has been stopped")
			return
		default:
			n, _, err := track.Read(buf)
			if err != nil {
				log.Println("Cannot read track:", err)
				return
			}

			packet := &rtp.Packet{}
			if err := packet.Unmarshal(buf[:n]); err != nil {
				log.Println("Failed to parse RTP:", err)
				return
			}

			writeErr := r.writer.Write(packet)

			if writeErr != nil {
				log.Println("Saving failed to OGG:", writeErr)
				return
			}
		}
	}
}
