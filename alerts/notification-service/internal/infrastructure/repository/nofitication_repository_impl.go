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

func castToObjectIds(ids ...string) ([]bson.ObjectID, error) {
	var objectIDs []bson.ObjectID
	var err error

	for _, id := range ids {
		uid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("Failed to convert id %s to ObjectID: %v (skipping)", id, err)
			continue
		}
		objectIDs = append(objectIDs, uid)
	}

	if len(objectIDs) == 0 {
		return []bson.ObjectID{}, fmt.Errorf("empty result list")
	}

	return objectIDs, err
}

// Delete implements repository.NotificationRepository.
func (r *notificationRepositoryImpl) Delete(ctx context.Context, clientID string, ids ...string) error {
	objectIDs, err := castToObjectIds(ids...)

	if err != nil {
		return err
	}

	filter := bson.M{
		"$and": []bson.M{
			{"_id": bson.M{"$in": objectIDs}},
			{"receiver._id": clientID},
		},
	}

	result, err := r.database.GetCollection().DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete notifications: %w", err)
	}

	log.Printf("Deleted %d for %s", result.DeletedCount, clientID)

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

func (r *notificationRepositoryImpl) MarkAsDelivered(ctx context.Context, ids ...string) error {
	if len(ids) == 0 {
		return fmt.Errorf("no ids provided")
	}

	var objectIDs []bson.ObjectID
	for _, id := range ids {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid hex string %q: %w", id, err)
		}
		objectIDs = append(objectIDs, oid)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{"$set": bson.M{"status": enums.DELIVERED}}

	res, err := r.database.GetCollection().UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update notifications: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no notifications found")
	}

	return nil
}

func (r *notificationRepositoryImpl) Save(ctx context.Context, n ...models.Notification) ([]*dto.NotificationDTO, error) {
	if len(n) == 0 {
		return nil, fmt.Errorf("no notifications to save")
	}

	res, err := r.database.GetCollection().InsertMany(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("failed to save notifications: %w", err)
	}

	if len(res.InsertedIDs) != len(n) {
		return nil, fmt.Errorf("number of inserted IDs does not match number of notifications")
	}

	for i, id := range res.InsertedIDs {
		oid, ok := id.(bson.ObjectID)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve inserted ID: unexpected type %T", id)
		}
		n[i].ID = oid
	}

	result := []*dto.NotificationDTO{}

	for _, notification := range n {
		result = append(result, notification.DTO())
	}

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
