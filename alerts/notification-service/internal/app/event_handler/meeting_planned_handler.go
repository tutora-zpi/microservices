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

type MeetingPlannedHandler struct {
	publisher interfaces.NotificationManager
	repo      repository.NotificationRepository
}

func NewMeetingPlannedHandler(publisher interfaces.NotificationManager,
	repo repository.NotificationRepository) interfaces.EventHandler {
	return &MeetingPlannedHandler{repo: repo, publisher: publisher}
}

func (c *MeetingPlannedHandler) Handle(ctx context.Context, body []byte) error {
	newEvent := meetinginvitation.PlannedMeetingEvent{}
	log.Printf("[%s] received\n", newEvent.Name())

	if err := json.Unmarshal(body, &newEvent); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	results, err := c.repo.Save(ctx, newEvent.Notifications()...)

	if err != nil || len(results) != 2 {
		log.Printf("An error occured during saving notification: %s\n", err.Error())
		return err
	}

	ids := []string{}
	for _, result := range results {
		if err := c.publisher.Push(*result); err != nil {
			return err
		}

		ids = append(ids, result.ID)
	}

	if err := c.repo.MarkAsDelivered(ctx, ids...); err != nil {
		return err
	}

	return nil
}
