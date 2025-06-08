package recorder

import (
	"log"
	"sync"
	"voice-service/internal/app/interfaces"
	"voice-service/internal/infrastructure/writer"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

type session struct {
	conn     *webrtc.PeerConnection
	writer   writer.Writer
	stopChan chan struct{}
	wg       sync.WaitGroup
}

type recorderImpl struct {
	cfg           *webrtc.Configuration
	writerFactory writer.WriterFactory

	sessions map[string]*session
	mutex    sync.Mutex
}

func NewRecorder(cfg *webrtc.Configuration, wf writer.WriterFactory) (interfaces.Recorder, error) {
	if cfg == nil {
		cfg = &webrtc.Configuration{}
	}

	return &recorderImpl{
		cfg:           cfg,
		writerFactory: wf,
		sessions:      make(map[string]*session),
	}, nil
}

func (r *recorderImpl) StartRecording(meetingID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.sessions[meetingID]; exists {
		return nil
	}

	peerConnection, err := webrtc.NewPeerConnection(*r.cfg)
	if err != nil {
		return err
	}

	stopChan := make(chan struct{})
	sess := &session{
		conn:     peerConnection,
		stopChan: stopChan,
	}

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		r.handleBuffer(track, meetingID)
	})

	r.sessions[meetingID] = sess
	return nil
}

func (r *recorderImpl) StopRecording(meetingID string) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	sess, exists := r.sessions[meetingID]
	if !exists {
		log.Println("No recording session found for meeting:", meetingID)
		return "", nil
	}

	close(sess.stopChan)
	sess.wg.Wait()

	err := sess.writer.Close()
	if err != nil {
		log.Println("Failed to close writer for meeting:", meetingID, err)
	}

	err = sess.conn.Close()
	if err != nil {
		log.Println("Failed to close peer connection for meeting:", meetingID, err)
	}

	path := r.sessions[meetingID].writer.GetPath()

	delete(r.sessions, meetingID)
	log.Println("Recording finished for meeting:", meetingID)
	return path, nil
}

func (r *recorderImpl) handleBuffer(track *webrtc.TrackRemote, meetingID string) {
	if track.Kind() != webrtc.RTPCodecTypeAudio {
		log.Println("Wrong track type:", track.Kind())
		return
	}

	r.mutex.Lock()
	sess, exists := r.sessions[meetingID]
	r.mutex.Unlock()

	if !exists {
		log.Println("No session found for meeting:", meetingID)
		return
	}

	if sess.writer == nil {
		writer, err := r.writerFactory(track, meetingID)
		if err != nil {
			log.Println("Failed to create writer:", err)
			return
		}
		sess.writer = writer
	}

	sess.wg.Add(1)
	go r.receiveAudio(track, sess)
}

func (r *recorderImpl) receiveAudio(track *webrtc.TrackRemote, sess *session) {
	defer sess.wg.Done()

	buf := make([]byte, 1400)
	for {
		select {
		case <-sess.stopChan:
			log.Println("Receiving audio stopped")
			return
		default:
			n, _, err := track.Read(buf)
			if err != nil {
				log.Println("Cannot read track:", err)
				return
			}

			packet := &rtp.Packet{}
			if err := packet.Unmarshal(buf[:n]); err != nil {
				log.Println("Failed to parse RTP packet:", err)
				return
			}

			if err := sess.writer.Write(packet); err != nil {
				log.Println("Failed to write packet:", err)
				return
			}
		}
	}
}
