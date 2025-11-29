package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	meetinginvitation "notification-serivce/internal/domain/event/meeting_invitation"
	"notification-serivce/internal/domain/repository"
)

type meetingEndedHandler struct {
	publisher interfaces.NotificationManager
	repo      repository.NotificationRepository
}

// Handle implements interfaces.EventHandler.
func (m *meetingEndedHandler) Handle(ctx context.Context, body []byte) error {
	newEvent := meetinginvitation.MeetingEndedEvent{}
	log.Printf("[%s] received\n", newEvent.Name())

	if err := json.Unmarshal(body, &newEvent); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	notifications := newEvent.EndedMeetingNotifications()

	for _, result := range notifications {
		if err := m.publisher.Push(*result.DTO()); err != nil {
			return err
		}

	}

	return nil
}

func NewMeetingEndedHandler(
	publisher interfaces.NotificationManager,
	repo repository.NotificationRepository,
) interfaces.EventHandler {
	return &meetingEndedHandler{repo: repo, publisher: publisher}
}
