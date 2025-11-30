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

	notificationForReceiver := *newEvent.NotificationForReceiver()
	nofiticationForSender := *newEvent.NotificationForSender()

	results, err := c.repo.Save(ctx, notificationForReceiver, nofiticationForSender)

	if err != nil {
		log.Printf("Failed to save notification: %v", err)
		return err
	}

	ids := []string{}
	for _, result := range results {
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
