package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/event/notes"
	"notification-serivce/internal/domain/repository"
)

type ResourceGeneratedHandler struct {
	publisher interfaces.NotificationManager
	repo      repository.NotificationRepository
}

func NewResourceGeneratedHandler(publisher interfaces.NotificationManager,
	repo repository.NotificationRepository) interfaces.EventHandler {
	return &ResourceGeneratedHandler{repo: repo, publisher: publisher}
}

func (c *ResourceGeneratedHandler) Handle(ctx context.Context, body []byte) error {
	var newEvent notes.ResourcesGeneratedEvent

	if err := json.Unmarshal(body, &newEvent); err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
		return fmt.Errorf("an error occured during casting into event struct")
	}

	ns, err := newEvent.Notifications()
	if err != nil {
		return err
	}

	results, err := c.repo.Save(ctx, ns...)

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
