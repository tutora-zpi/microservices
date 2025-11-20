package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"recorder-service/internal/app/interfaces"
	"recorder-service/internal/app/interfaces/handler"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/broker"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/recorder"
	"recorder-service/internal/domain/repository"
	"recorder-service/internal/infrastructure/ffmpeg"
	"recorder-service/internal/infrastructure/s3"
	"sync"
)

type stopRecordingMeetingHandler struct {
	botService service.BotService
	repo       repository.VoiceSessionMetadataRepository
	s3service  s3.S3Service
	broker     interfaces.Broker

	exchange string
}

// Handle implements interfaces.EventHandler.
func (s *stopRecordingMeetingHandler) Handle(ctx context.Context, body []byte) error {
	var evt event.StopRecordingMeetingEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode: %v", err)
		return fmt.Errorf("failed to decode: %s", evt.Name())
	}

	infos, err := s.botService.RemoveBot(ctx, evt.RoomID)
	if err != nil || len(infos) < 1 {
		log.Printf("%v", err)
	}

	outputPath, err := s.MergeRecordings(infos)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("Recordings merged and saved: %s", outputPath)

	var wg sync.WaitGroup
	var errors chan error = make(chan error, 2)

	wg.Go(func() {
		basePath := infos[0].BasePath
		keys, err := s.s3service.PutObject(ctx, basePath)
		if err != nil {
			errors <- err
		}

		evt := event.NewRecordingsUploaded(keys)
		log.Printf("RecordingsUploaded body: %v", *evt)

		dest := broker.NewExchangeDestination(evt, s.exchange)
		err = s.broker.Publish(ctx, evt, dest)
		if err != nil {
			errors <- err
		}
	})

	wg.Go(func() {
		var err = s.repo.AppendAudioName(ctx, evt.RoomID, outputPath)
		if err != nil {
			errors <- err
		}
	})

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *stopRecordingMeetingHandler) MergeRecordings(infos []recorder.RecordingInfo) (string, error) {
	if err := ffmpeg.AddSilence(infos); err != nil {
		return "", err
	}

	return ffmpeg.MixAudio(infos)
}

func NewStopRecordingMeetingHandler(
	botService service.BotService,
	repo repository.VoiceSessionMetadataRepository,
	s3service s3.S3Service,
	broker interfaces.Broker,
	exchange string,
) handler.EventHandler {
	return &stopRecordingMeetingHandler{botService: botService, repo: repo, s3service: s3service, broker: broker, exchange: exchange}
}
