package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"recorder-service/internal/app/interfaces/handler"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/repository"
	"sync"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type stopRecordingMeetingHandler struct {
	botService service.BotService
	repo       repository.VoiceSessionMetadataRepository
}

// Handle implements interfaces.EventHandler.
func (s *stopRecordingMeetingHandler) Handle(ctx context.Context, body []byte) error {
	var evt event.StopRecordingMeetingEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode: %v", err)
		return fmt.Errorf("failed to decode: %s", evt.Name())
	}

	bot, ok := s.botService.GetBot(evt.RoomID)
	if !ok {
		return fmt.Errorf("no bot in room %s", evt.RoomID)
	}

	paths, err := bot.Recorder().StopRecording(ctx, evt.RoomID)
	if err != nil || len(paths) < 2 {
		log.Printf("Failed to stop recording: %v", err)
		return fmt.Errorf("failed to stop recording")
	}

	fileName := s.GetFileName(paths[0])
	fileDest := path.Join("recordings", evt.RoomID, fileName)
	var wg sync.WaitGroup

	wg.Go(func() {
		if err := s.MergeRecordings(paths, fileDest); err != nil {
			log.Printf("Failed to update merge recordings: %v", err)
		}
	})

	wg.Go(func() {
		if err := s.repo.AppendAudioName(ctx, evt.RoomID, fileName); err != nil {
			log.Printf("Failed to update metadata: %v", err)
		}
	})

	wg.Wait()

	return nil
}

func (s *stopRecordingMeetingHandler) GetFileName(examplePath string) string {
	ext := path.Ext(examplePath)

	now := time.Now().UTC().UnixNano()
	fileName := fmt.Sprintf("merged_%d%s", now, ext)

	return fileName
}

func (s *stopRecordingMeetingHandler) MergeRecordings(paths []string, fileDest string) error {
	var stream *ffmpeg.Stream
	for _, path := range paths {
		stream = ffmpeg.Input(path)
	}

	err := stream.Filter("amix", nil, ffmpeg.KwArgs{"inputs": len(paths), "duration": "longest"}).
		Output(fileDest).
		Run()

	if err != nil {
		log.Printf("Failed to merge audio: %v", err)
		return fmt.Errorf("failed to run mix")
	}

	return nil
}

func NewStopRecordingMeetingHandler(botService service.BotService, repo repository.VoiceSessionMetadataRepository) handler.EventHandler {
	return &stopRecordingMeetingHandler{botService: botService, repo: repo}
}
