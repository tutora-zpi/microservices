package usecase

import (
	"fmt"
	"log"
	"recorder-service/internal/app/interfaces"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/repository"
)

// react on a event
type CreateMeetingSession struct {
	repo repository.VoiceMeetingRepository
}

// Exec implements interfaces.UseCaseHandler.
func (c *CreateMeetingSession) Exec(body []byte) error {
	var ev *event.EventWrapper
	ev = ev.FromJson(body)

	var dest event.MeetingStartedEvent
	if err := ev.DecodeBody(dest); err != nil {
		log.Printf("Failed to decode payload %s\n", err.Error())
		return fmt.Errorf("failed to decode payload")
	}

	_, err := c.repo.CreateSession(dest)

	if err != nil {
		log.Println("Failed to create session")
		return err
	}

	return nil
}

func NewCreateMeetingSession(repo repository.VoiceMeetingRepository) interfaces.UseCaseHandler {
	return &CreateMeetingSession{
		repo: repo,
	}
}
