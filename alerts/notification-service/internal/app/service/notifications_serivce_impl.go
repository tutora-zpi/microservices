package service

import (
	"context"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/repository"
)

type notificationSerivceImpl struct {
	repo repository.NotificationRepository
}

// DeleteNotifications implements interfaces.NotificationSerivce.
func (n *notificationSerivceImpl) DeleteNotifications(ctx context.Context, req *dto.DeleteNotificationsDTO, clientID string) error {
	log.Println("Deleting notifications...")

	return n.repo.Delete(ctx, clientID, req.IDs...)
}

// FetchNotifications implements interfaces.NotificationSerivce.
func (n *notificationSerivceImpl) FetchNotifications(ctx context.Context, req *dto.FetchNotificationsDTO) ([]dto.NotificationDTO, error) {
	log.Println("Fetching notifications...")

	return n.repo.Get(ctx, req.ReceiverID, req.LastNotificationID, req.Limit)
}

func NewNotificationSerivce(repo repository.NotificationRepository) interfaces.NotificationSerivce {
	return &notificationSerivceImpl{repo: repo}
}
