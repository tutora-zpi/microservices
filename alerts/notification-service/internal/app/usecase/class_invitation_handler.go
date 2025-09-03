package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/event"
	"notification-serivce/internal/domain/repository"
)

type ClassInvitationHandler struct {
	publisher interfaces.NotificationManager
	repo      repository.NotificationRepository
}

func (c *ClassInvitationHandler) Handle(body []byte) error {
	ctx := context.Background()
	event := event.ClassInvitationEvent{}
	var err error

	if err = json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	log.Println("Successfully handled")

	var dto *dto.NotificationDTO

	dto, err = c.repo.Save(ctx, event.Notification())
	if err != nil {
		return err
	}

	log.Println("DTO", *dto)

	if err = c.publisher.Push(*dto); err != nil {
		return err
	}

	if err = c.repo.MarkAsDelivered(ctx, dto.ID); err != nil {
		return err
	}

	return nil
}

func NewClassInvitationHandler(publisher interfaces.NotificationManager, repo repository.NotificationRepository) interfaces.EventHandler {
	return &ClassInvitationHandler{
		publisher: publisher,
		repo:      repo,
	}
}
