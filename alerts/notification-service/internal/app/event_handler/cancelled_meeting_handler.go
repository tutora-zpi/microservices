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

type CancelledMeetingHandler struct {
	publisher interfaces.NotificationManager
	repo      repository.NotificationRepository
}

func NewCancelledMeetingHandler(publisher interfaces.NotificationManager,
	repo repository.NotificationRepository) interfaces.EventHandler {
	return &CancelledMeetingHandler{repo: repo, publisher: publisher}
}

func (c *CancelledMeetingHandler) Handle(ctx context.Context, body []byte) error {
	newEvent := meetinginvitation.CancelledMeetingEvent{}
	log.Printf("[%s] received\n", newEvent.Name())

	if err := json.Unmarshal(body, &newEvent); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	notifications := newEvent.Notifications()

	results, err := c.repo.Save(ctx, notifications...)

	if err != nil {
		log.Printf("Failed to save notification: %v", err)
		return err
	}

	ids := []string{}
	for _, result := range results {
		log.Printf("Notfi: %v", *result)
		if err := c.publisher.Push(*result); err != nil {
			log.Printf("Failed to push notification to user: %v", err)
			continue
		}

		ids = append(ids, result.ID)
	}

	if err := c.repo.MarkAsDelivered(ctx, ids...); err != nil {
		return err
	}

	return nil
}
