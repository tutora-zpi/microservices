package service

import (
	"context"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/repository"
	"time"
)

type notificationSerivceImpl struct {
	repo repository.NotificationRepository
}

// DeleteNotifications implements interfaces.NotificationSerivce.
func (n *notificationSerivceImpl) DeleteNotifications(req *dto.DeleteNotificationsDTO, clientID string) error {
	log.Println("Deleting notifications...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return n.repo.Delete(ctx, clientID, req.IDs...)
}

// FetchNotifications implements interfaces.NotificationSerivce.
func (n *notificationSerivceImpl) FetchNotifications(req *dto.FetchNotificationsDTO) ([]dto.NotificationDTO, error) {
	log.Println("Fetching notifications...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return n.repo.Get(ctx, req.ReceiverID, req.LastNotificationID, req.Limit)
}

func NewNotificationSerivce(repo repository.NotificationRepository) interfaces.NotificationSerivce {
	return &notificationSerivceImpl{repo: repo}
}
