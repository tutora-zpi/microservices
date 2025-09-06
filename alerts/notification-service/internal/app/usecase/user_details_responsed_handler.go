package usecase

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

type UserDetailsResponsedHandler struct {
	broker interfaces.Broker
	repo   repository.NotificationRepository
}

func NewUserDetailsResponsedHandler(broker interfaces.Broker,
	repo repository.NotificationRepository) interfaces.EventHandler {
	return &UserDetailsResponsedHandler{repo: repo, broker: broker}
}

func (u *UserDetailsResponsedHandler) Handle(body []byte) error {
	log.Println("UserDetailsResponsedHandler received")
	ctx := context.Background()

	newEvent := classinvitation.UserDetailsRespondedEvent{}
	if err := json.Unmarshal(body, &newEvent); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	fieldsToUpdate := newEvent.FieldsToUpdate()

	result, err := u.repo.Update(ctx, fieldsToUpdate, newEvent.ID)
	if err != nil {
		log.Printf("An error occured during saving notification: %s", err.Error())
		return err
	}

	readyNotifcation := classinvitation.NewClassInvitationReadyEvent(result)

	if err := u.broker.Publish(event.NewEventWrapper(readyNotifcation)); err != nil {
		log.Printf("An error occured during publishing %s: %s\n", readyNotifcation.Name(), err.Error())
		return err
	}

	return nil
}
