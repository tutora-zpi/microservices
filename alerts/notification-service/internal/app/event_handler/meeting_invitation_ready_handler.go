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

type MeetingInvitationReadyEventHandler struct {
	publisher interfaces.NotificationManager
	repo      repository.NotificationRepository
}

func NewMeetingInvitationReadyEventHandler(publisher interfaces.NotificationManager,
	repo repository.NotificationRepository) interfaces.EventHandler {
	return &MeetingInvitationReadyEventHandler{repo: repo, publisher: publisher}
}

func (m *MeetingInvitationReadyEventHandler) Handle(ctx context.Context, body []byte) error {
	newEvent := meetinginvitation.MeetingStartedEvent{}
	log.Printf("[%s] received", newEvent.Name())

	if err := json.Unmarshal(body, &newEvent); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	notifications := newEvent.Notifications()

	results, err := m.repo.Save(ctx, notifications...)
	if err != nil {
		log.Printf("An error occured during saving partial notification: %s\n", err.Error())
		return err
	}

	ids := []string{}
	for _, result := range results {
		if err := m.publisher.Push(*result); err != nil {
			return err
		}

		ids = append(ids, result.ID)
	}

	if err := m.repo.MarkAsDelivered(ctx, ids...); err != nil {
		return err
	}

	return nil
}
