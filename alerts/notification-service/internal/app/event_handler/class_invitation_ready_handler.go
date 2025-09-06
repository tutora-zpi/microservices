package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/dto"
	classinvitation "notification-serivce/internal/domain/event/class_invitation"
	"notification-serivce/internal/domain/repository"
)

type ClassInvitationReadyHandler struct {
	publisher interfaces.NotificationManager
	repo      repository.NotificationRepository
}

func (c *ClassInvitationReadyHandler) Handle(body []byte) error {
	ctx := context.Background()
	event := classinvitation.ClassInvitationReadyEvent{}
	var err error

	if err = json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	log.Println("Successfully handled")

	var dto *dto.NotificationDTO = event.GetReadyNotification()

	if err = c.publisher.Push(*dto); err != nil {
		return err
	}

	if err = c.repo.MarkAsDelivered(ctx, dto.ID); err != nil {
		return err
	}

	return nil
}

func NewClassInvitationReadyHandler(publisher interfaces.NotificationManager, repo repository.NotificationRepository) interfaces.EventHandler {
	return &ClassInvitationReadyHandler{
		publisher: publisher,
		repo:      repo,
	}
}
