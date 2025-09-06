package service

import (
	"context"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/repository"
	"notification-serivce/internal/domain/requests"
)

type notificationSerivceImpl struct {
	repo repository.NotificationRepository
}

// DeleteNotifications implements interfaces.NotificationSerivce.
func (n *notificationSerivceImpl) DeleteNotifications(req *requests.DeleteNotificationsRequest, clientID string) error {
	log.Println("Deleting notifications...")
	ctx := context.Background()
	return n.repo.Delete(ctx, clientID, req.IDs...)
}

// FetchNotifications implements interfaces.NotificationSerivce.
func (n *notificationSerivceImpl) FetchNotifications(req *requests.FetchNotificationsRequest) ([]dto.NotificationDTO, error) {
	log.Println("Fetching notifications...")
	ctx := context.Background()

	return n.repo.Get(ctx, req.ReceiverID, req.LastNotificationID, req.Limit)
}

func NewNotificationSerivce(repo repository.NotificationRepository) interfaces.NotificationSerivce {
	return &notificationSerivceImpl{repo: repo}
}
