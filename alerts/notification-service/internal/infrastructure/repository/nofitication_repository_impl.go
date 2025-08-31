package repository

import (
	"context"
	"fmt"
	"log"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/models"
	"notification-serivce/internal/domain/repository"
	"notification-serivce/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type notificationRepositoryImpl struct {
	database *database.Database
}

// MarkAsDelivered implements repository.NotificationRepository.
func (r *notificationRepositoryImpl) MarkAsDelivered(ctx context.Context, id string) error {
	uid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid hex string")
	}

	filter := bson.M{"_id": uid}
	update := bson.M{"$set": bson.M{"status": enums.Delivered}}

	res := r.database.GetCollection().FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("notification not found")
		}
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

func NewNotificationRepository(db *database.Database) repository.NotificationRepository {
	return &notificationRepositoryImpl{
		database: db,
	}
}

// Save implements repository.NotificationRepository.
func (this *notificationRepositoryImpl) Save(ctx context.Context, n *models.Notification) (*dto.NotificationDTO, error) {
	res, err := this.database.GetCollection().InsertOne(ctx, n)
	if err != nil {
		log.Printf("Failed to save notification: %s\n", err.Error())
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	if uid, ok := res.InsertedID.(bson.ObjectID); ok {
		n.ID = uid
	} else {
		log.Printf("Unexpected InsertedID type: %T", res.InsertedID)
		return nil, fmt.Errorf("failed to retrieve inserted ID: unexpected type %T", res.InsertedID)
	}

	result := n.DTO()
	return &result, nil
}
