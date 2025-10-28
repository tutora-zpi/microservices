package recorder

import (
	"context"
	"fmt"
	"recorder-service/internal/domain/client"
	"recorder-service/internal/domain/recorder"
	"recorder-service/internal/infrastructure/webrtc/writer"
	"sync"

	"github.com/pion/webrtc/v3"
)

type TrackInfo struct {
	track  *webrtc.TrackRemote
	cancel context.CancelFunc
	writer writer.Writer
}

type RoomConnections struct {
	clients   map[string]client.Client
	tracks    map[string]*TrackInfo
	newTracks chan *TrackInfo
}

type RecorderClient struct {
	rooms map[string]*RoomConnections // roomID -> RoomConnections
	mutex sync.Mutex
	wg    sync.WaitGroup
}

// RegisterNewRoom implements recorder.Recorder.
func (r *RecorderClient) RegisterNewRoom(ctx context.Context, roomID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.rooms[roomID]; !ok {
		r.rooms[roomID] = &RoomConnections{
			clients:   make(map[string]client.Client),
			tracks:    make(map[string]*TrackInfo),
			newTracks: make(chan *TrackInfo, 100),
		}
	}

	return nil
}

// StartRecording implements recorder.Recorder.
func (r *RecorderClient) StartRecording(ctx context.Context, roomID string) error {
	r.mutex.Lock()
	roomConn, ok := r.rooms[roomID]
	if !ok {
		r.mutex.Unlock()
		return fmt.Errorf("room: %s does not exist", roomID)
	}
	r.mutex.Unlock()

	r.wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case trackInfo := <-roomConn.newTracks:
				if trackInfo == nil || trackInfo.track == nil {
					continue
				}
				ctxRec, cancel := context.WithCancel(ctx)
				trackInfo.cancel = cancel
				go trackInfo.writer.Write(ctxRec, trackInfo.track)
			}
		}
	})

	r.mutex.Lock()
	existingTracks := make([]*TrackInfo, 0, len(roomConn.tracks))
	for _, t := range roomConn.tracks {
		existingTracks = append(existingTracks, t)
	}
	r.mutex.Unlock()

	for _, trackInfo := range existingTracks {
		if trackInfo == nil || trackInfo.track == nil {
			continue
		}
		ctxRec, cancel := context.WithCancel(ctx)
		trackInfo.cancel = cancel
		go trackInfo.writer.Write(ctxRec, trackInfo.track)
	}

	return nil
}

// AddRecordedClient implements recorder.Recorder.
func (r *RecorderClient) AddRecordedClient(ctx context.Context, writerFactory writer.WriterFactory, roomID, userID string, client client.Client) error {
	r.mutex.Lock()
	roomConn, ok := r.rooms[roomID]
	if !ok {
		r.mutex.Unlock()
		return fmt.Errorf("room %s does not exist", roomID)
	}

	if _, ok := roomConn.clients[userID]; ok {
		r.mutex.Unlock()
		return nil
	}

	roomConn.clients[userID] = client
	r.mutex.Unlock()

	client.OnTrack(func(track *webrtc.TrackRemote) {
		writer, _ := writerFactory(roomID, userID, track)
		trackInfo := &TrackInfo{
			track:  track,
			writer: writer,
		}

		r.mutex.Lock()
		roomConn.tracks[userID] = trackInfo
		r.mutex.Unlock()

		roomConn.newTracks <- trackInfo
	})

	return nil
}

// StopRecording implements recorder.Recorder.
func (r *RecorderClient) StopRecording(ctx context.Context, roomID string) ([]string, error) {
	r.mutex.Lock()
	roomConn, ok := r.rooms[roomID]
	if !ok {
		r.mutex.Unlock()
		return nil, fmt.Errorf("room: %s does not exist", roomID)
	}
	r.mutex.Unlock()

	paths := []string{}
	for _, trackInfo := range roomConn.tracks {
		if trackInfo != nil && trackInfo.cancel != nil {
			trackInfo.cancel()
		}
		paths = append(paths, trackInfo.writer.GetPath())
	}

	r.mutex.Lock()
	delete(r.rooms, roomID)
	r.mutex.Unlock()

	return paths, nil
}

func NewRecorderClient() recorder.Recorder {
	return &RecorderClient{
		rooms: make(map[string]*RoomConnections),
	}
}
