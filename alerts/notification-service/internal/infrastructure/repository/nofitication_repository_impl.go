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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type notificationRepositoryImpl struct {
	database *database.Database
}

// Delete implements repository.NotificationRepository.
func (r *notificationRepositoryImpl) Delete(ctx context.Context, clientID string, ids ...string) error {
	var objectIDs []bson.ObjectID
	var err error
	var uid bson.ObjectID

	for _, id := range ids {
		uid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("Failed to convert id %s to ObjectID: %v (skipping)", id, err)
			continue
		}
		objectIDs = append(objectIDs, uid)
	}

	if len(objectIDs) == 0 {
		return nil
	}

	uid, err = bson.ObjectIDFromHex(clientID)
	if err != nil {
		return fmt.Errorf("invalid hex string")
	}

	filter := bson.M{
		"$and": []bson.M{
			{"_id": bson.M{"$in": objectIDs}},
			{"receiver._id": uid},
		},
	}

	_, err = r.database.GetCollection().DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete notifications: %w", err)
	}

	return nil
}

func (r *notificationRepositoryImpl) Update(ctx context.Context, fields map[string]any, id string) (*dto.NotificationDTO, error) {
	uid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string")
	}

	filter := bson.M{"_id": uid}

	update := bson.M{"$set": fields}

	res := r.database.GetCollection().FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	var result models.Notification
	if err := res.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode notification: %w", err)
	}

	return result.DTO(), nil
}

func NewNotificationRepository(db *database.Database) repository.NotificationRepository {
	return &notificationRepositoryImpl{
		database: db,
	}
}

func (r *notificationRepositoryImpl) MarkAsDelivered(ctx context.Context, id string) error {
	uid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid hex string")
	}

	filter := bson.M{"_id": uid}
	update := bson.M{"$set": bson.M{"status": enums.DELIVERED}}

	res := r.database.GetCollection().FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("notification not found")
		}
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

func (this *notificationRepositoryImpl) Save(ctx context.Context, n *models.Notification) (*dto.NotificationDTO, error) {
	res, err := this.database.GetCollection().InsertOne(ctx, n)
	if err != nil {
		log.Printf("Failed to save notification: %s\n", err.Error())
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	if uid, ok := res.InsertedID.(bson.ObjectID); ok {
		n.ID = uid
	} else {
		return nil, fmt.Errorf("failed to retrieve inserted ID: unexpected type %T", res.InsertedID)
	}

	result := n.DTO()
	return result, nil
}

func (r *notificationRepositoryImpl) Get(ctx context.Context, receiverID string, lastNotificationID *string, limit int) ([]dto.NotificationDTO, error) {
	filter := bson.M{"receiver._id": receiverID}

	if lastNotificationID != nil {
		uid, err := bson.ObjectIDFromHex(*lastNotificationID)
		if err != nil {
			return nil, fmt.Errorf("invalid hex string")
		}
		filter["_id"] = bson.M{"$lt": uid}

	}

	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.database.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer cursor.Close(ctx)

	decoded, err := r.decodeNotifications(ctx, cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode notifications: %w", err)
	}

	return decoded, nil
}

func (r *notificationRepositoryImpl) decodeNotifications(ctx context.Context, cursor *mongo.Cursor) ([]dto.NotificationDTO, error) {
	var dtos []dto.NotificationDTO

	for cursor.Next(ctx) {
		var result models.Notification
		if err := cursor.Decode(&result); err != nil {
			return dtos, fmt.Errorf("failed to decode notification: %w", err)
		}
		dtos = append(dtos, *result.DTO())
	}

	if err := cursor.Err(); err != nil {
		return dtos, fmt.Errorf("cursor error: %w", err)
	}

	return dtos, nil
}
