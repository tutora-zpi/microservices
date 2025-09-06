package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/event"
	classinvitation "notification-serivce/internal/domain/event/class_invitation"
	"notification-serivce/internal/domain/repository"
)

type ClassInvitationCreatedHandler struct {
	broker interfaces.Broker
	repo   repository.NotificationRepository
}

func NewClassInvitationCreatedHandler(broker interfaces.Broker,
	repo repository.NotificationRepository) interfaces.EventHandler {
	return &ClassInvitationCreatedHandler{repo: repo, broker: broker}
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
	if err != nil {
		log.Printf("An error occured during saving partial notification: %s\n", err.Error())
		return err
	}

	userDetails := classinvitation.NewUserDetailsRequestedEvent(result)

	wrapped := event.NewEventWrapper(userDetails)

	if err := c.broker.Publish(wrapped); err != nil {
		log.Printf("An error occured during publishing %s: %s\n", userDetails.Name(), err.Error())
		return err
	}

	return nil
}
