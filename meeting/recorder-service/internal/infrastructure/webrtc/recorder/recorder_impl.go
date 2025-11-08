package recorder

import (
	"context"
	"log"
	"recorder-service/internal/domain/recorder"
	"recorder-service/internal/infrastructure/webrtc/writer"
	"sync"

	"github.com/pion/webrtc/v3"
)

type UserID = string

type recorderImpl struct {
	mu               sync.Mutex
	recordingDetails map[UserID]*recorder.Detail
}

// StartRecording implements recorder.Recorder.
func (r *recorderImpl) StartRecording(ctx context.Context, roomID, userID string, track *webrtc.TrackRemote, newWriter writer.WriterFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.recordingDetails == nil {
		r.recordingDetails = make(map[UserID]*recorder.Detail)
	}

	if _, ok := r.recordingDetails[userID]; ok {
		log.Printf("Recorder: writer for user %s already exists", userID)
		return
	}

	w, err := newWriter(roomID)
	if err != nil {
		log.Printf("Recorder: failed to create writer for user %s: %v", userID, err)
		return
	}

	userCtx, cancel := context.WithCancel(ctx)
	detail := recorder.NewDetail(w, cancel)
	detail.SetJoinTime()

	r.recordingDetails[userID] = detail
	if err := w.Write(userCtx, userID, track); err != nil {
		log.Printf("An error occurred during writing: %v", err)
	}
}

// StopRecording implements recorder.Recorder.
func (r *recorderImpl) StopRecording(userID string) *recorder.RecordingInfo {
	r.mu.Lock()
	defer r.mu.Unlock()

	if detail, ok := r.recordingDetails[userID]; ok {
		detail.Cancel()
		detail.Writer.Close()
		detail.SetLeftTime()

		basePath := detail.Writer.GetPath()

		info := recorder.NewRecordingInfo(detail, basePath, userID, detail.Writer.GetExt())

		delete(r.recordingDetails, userID)

		log.Printf("Stopped recording user %s, returning info", userID)
		return &info
	}

	return nil
}

// StopAll implements recorder.Recorder.
func (r *recorderImpl) StopAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, detail := range r.recordingDetails {
		detail.Cancel()
		detail.Writer.Close()
	}

	r.recordingDetails = make(map[UserID]*recorder.Detail)
}

func NewRecorderClient() recorder.Recorder {
	return &recorderImpl{
		recordingDetails: make(map[UserID]*recorder.Detail),
	}
}
