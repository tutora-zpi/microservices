package repoimpl

import (
	"context"
	"fmt"
	"log"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/model"
	"recorder-service/internal/domain/repository"
	"recorder-service/internal/infrastructure/mongo"
	"recorder-service/pkg"
	"sort"

	"go.mongodb.org/mongo-driver/v2/bson"
	mongodb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type voiceMeetingRepoImpl struct {
	collection *mongodb.Collection
}

// AppendAudioName implements repository.VoiceSessionMetadataRepository.
func (v *voiceMeetingRepoImpl) AppendAudioName(ctx context.Context, meetingID string, audioName string) (*dto.VoiceSessionMetadataDTO, error) {
	name := pkg.GetFileName(audioName)

	filter := bson.M{"meetingId": meetingID}
	update := bson.M{"$set": bson.M{"audioName": name}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var model model.VoiceSessionMetadata

	err := v.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to update metadata with meetingId %s: %w", meetingID, err)
	}

	return model.DTO(), nil
}

// CreateSessionMetadata implements repository.VoiceSessionMetadataRepository.
func (v *voiceMeetingRepoImpl) CreateSessionMetadata(ctx context.Context, event event.MeetingStartedEvent) (*dto.VoiceSessionMetadataDTO, error) {
	model := model.NewVoiceSession(event)

	res, err := v.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("failed to save metadata")
	}

	oid, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("unexpected inserted id type")
	}

	model.ID = oid

	return model.DTO(), nil
}

// FetchSessionMetadata implements repository.VoiceSessionMetadataRepository.
func (v *voiceMeetingRepoImpl) FetchSessionMetadata(ctx context.Context, meetingID string, limit int64, lastFetchedID *string) ([]*dto.VoiceSessionMetadataDTO, error) {
	filter := bson.M{"meetingId": meetingID}

	if lastFetchedID != nil {
		objID, err := bson.ObjectIDFromHex(*lastFetchedID)
		if err != nil {
			return nil, fmt.Errorf("invalid format of string id: %w", err)
		}

		filter["_id"] = bson.M{"$lt": objID}
	}

	opts := options.Find().
		SetSort(bson.M{"_id": -1}).
		SetLimit(limit)

	cursor, err := v.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Failed to find: %v", err)
		return nil, fmt.Errorf("not found")
	}

	var models []model.VoiceSessionMetadata

	err = cursor.All(ctx, &models)
	if err != nil {
		log.Printf("Failed to decode data: %v", err)
		return nil, fmt.Errorf("failed to decode")
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("not found")
	}

	var dtos []*dto.VoiceSessionMetadataDTO = make([]*dto.VoiceSessionMetadataDTO, len(models))
	for i, model := range models {
		dtos[i] = model.DTO()
	}

	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].StartedAt.Before(*dtos[j].EndedAt)
	})

	return dtos, nil
}

func NewVoiceMeetingRepository(client *mongodb.Client, mongoConfig mongo.MongoConfig) repository.VoiceSessionMetadataRepository {
	collection := client.Database(mongoConfig.Database).Collection(mongoConfig.SessionMetadataCollection)

	return &voiceMeetingRepoImpl{collection: collection}
}
