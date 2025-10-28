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

	"go.mongodb.org/mongo-driver/v2/bson"
	mongodb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type voiceMeetingRepoImpl struct {
	collection *mongodb.Collection
}

// AppendAudioName implements repository.VoiceSessionMetadataRepository.
func (v *voiceMeetingRepoImpl) AppendAudioName(ctx context.Context, id string, audioName string) error {
	res, err := v.collection.UpdateByID(ctx, id, bson.M{"$set": bson.M{"audioName": audioName}})
	if err != nil || res.ModifiedCount != 1 {
		log.Printf("Error occurred during appending audio name")
		return fmt.Errorf("failed to update metadata with id: %s", id)
	}

	return nil
}

// CreateSessionMetadata implements repository.VoiceSessionMetadataRepository.
func (v *voiceMeetingRepoImpl) CreateSessionMetadata(ctx context.Context, event event.MeetingStartedEvent) (*dto.VoiceSessionMetadataDTO, error) {
	newMetadata := model.NewVoiceSession(event)

	_, err := v.collection.InsertOne(ctx, newMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to save metadata")
	}

	return newMetadata.DTO(), nil
}

// FetchSessionMetadata implements repository.VoiceSessionMetadataRepository.
func (v *voiceMeetingRepoImpl) FetchSessionMetadata(ctx context.Context, classID string, limit int64, lastFetchedMeetingID *string) ([]*dto.VoiceSessionMetadataDTO, error) {
	filter := bson.M{"classId": classID}

	if lastFetchedMeetingID != nil {
		filter["meetingId"] = lastFetchedMeetingID
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

	return dtos, nil
}

func NewVoiceMeetingRepository(client *mongodb.Client, mongoConfig mongo.MongoConfig) repository.VoiceSessionMetadataRepository {
	collection := client.Database(mongoConfig.Database).Collection(mongoConfig.SessionMetadataCollection)

	return &voiceMeetingRepoImpl{collection: collection}
}
