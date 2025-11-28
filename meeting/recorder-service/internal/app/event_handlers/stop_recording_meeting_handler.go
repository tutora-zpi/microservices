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
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/recorder"
	"recorder-service/internal/domain/repository"
	"recorder-service/internal/infrastructure/ffmpeg"
	"recorder-service/internal/infrastructure/s3"
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
	var err error

	if err = json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode: %v", err)
		return fmt.Errorf("failed to decode: %s", evt.Name())
	}

	infos, err := s.botService.RemoveBot(ctx, evt.RoomID)
	if err != nil {
		log.Printf("%v", err)
	}

	if len(infos) < 1 {
		return fmt.Errorf("no info about tracks")
	}

	outputPath, err := s.MergeRecordings(infos)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("Recordings merged and saved: %s", outputPath)

	basePath := infos[0].BasePath
	keys, err := s.s3service.PutObject(ctx, basePath)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	var updated *dto.VoiceSessionMetadataDTO
	updated, err = s.repo.AppendAudioName(ctx, evt.RoomID, outputPath)
	if err != nil {
		return fmt.Errorf("failed to update audio name: %w", err)
	}

	uploadEvent := event.NewRecordingsUploaded(keys, updated.ClassID, updated.MeetingID)
	log.Printf("RecordingsUploaded body: %v", *uploadEvent)

	dest := broker.NewExchangeDestination(uploadEvent, s.exchange)
	err = s.broker.Publish(ctx, uploadEvent, dest)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
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
