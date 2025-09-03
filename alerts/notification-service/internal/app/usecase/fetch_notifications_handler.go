package usecase

import (
	"context"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/query"
	"notification-serivce/internal/domain/repository"
)

type FetchNotificationsHandler struct {
	repo repository.NotificationRepository
}

// Execute implements interfaces.QueryHandler.
func (f *FetchNotificationsHandler) Execute(q any) (any, error) {
	log.Println("Executing query: FetchNotifications")
	ctx := context.Background()

	body, ok := q.(*query.FetchNotificationsQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type: expected *FetchNotificationsQuery")
	}

	return f.repo.Get(ctx, body.ReceiverID, body.LastNotificationID, body.Limit)
}

func NewFetchNotificationsHandler(repo repository.NotificationRepository) interfaces.QueryHandler {
	return &FetchNotificationsHandler{
		repo: repo,
	}
}
