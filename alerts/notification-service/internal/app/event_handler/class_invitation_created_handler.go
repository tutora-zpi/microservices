package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	classinvitation "notification-serivce/internal/domain/event/class_invitation"
	"notification-serivce/internal/domain/repository"
)

type ClassInvitationCreatedHandler struct {
	publisher interfaces.NotificationManager
	repo      repository.NotificationRepository
}

func NewClassInvitationCreatedHandler(publisher interfaces.NotificationManager,
	repo repository.NotificationRepository) interfaces.EventHandler {
	return &ClassInvitationCreatedHandler{repo: repo, publisher: publisher}
}

func (c *ClassInvitationCreatedHandler) Handle(ctx context.Context, body []byte) error {
	newEvent := classinvitation.ClassInvitationCreatedEvent{}
	log.Printf("[%s] received\n", newEvent.Name())

	if err := json.Unmarshal(body, &newEvent); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	results, err := c.repo.Save(ctx, *newEvent.NotificationForReceiver(), *newEvent.NotificationForSender())

	if err != nil || len(results) != 2 {
		log.Printf("An error occured during saving notification: %s\n", err.Error())
		return err
	}

	for _, result := range results {
		if err = c.publisher.Push(*result); err != nil {
			return err
		}

		if err = c.repo.MarkAsDelivered(ctx, result.ID); err != nil {
			return err
		}
	}

	return nil
}
