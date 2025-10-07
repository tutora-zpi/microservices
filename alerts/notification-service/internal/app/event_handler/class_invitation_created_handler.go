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

func (c *ClassInvitationCreatedHandler) Handle(body []byte) error {
	log.Println("ClassInvitationCreated received")
	ctx := context.Background()

	newEvent := classinvitation.ClassInvitationCreatedEvent{}

	if err := json.Unmarshal(body, &newEvent); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	result, err := c.repo.Save(ctx, newEvent.Notification())
	dto := result[0]
	if err != nil || dto == nil {
		log.Printf("An error occured during saving partial notification: %s\n", err.Error())
		return err
	}

	if err = c.publisher.Push(*dto); err != nil {
		return err
	}

	if err = c.repo.MarkAsDelivered(ctx, dto.ID); err != nil {
		return err
	}

	return nil
}
