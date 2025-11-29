package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/models"
	"meeting-scheduler-service/internal/domain/repository"
	"meeting-scheduler-service/internal/infrastructure/cache"
	"time"

	"github.com/redis/go-redis/v9"
)

type meetingRepoImpl struct {
	client *redis.Client
}

// Exists implements repository.MeetingRepository.
func (m *meetingRepoImpl) Exists(ctx context.Context, classID string) bool {
	key := cache.MeetingKey(classID)

	foundNumber, _ := m.client.Exists(ctx, key).Result()
	return foundNumber > 0
}

// Append implements repository.MeetingRepository.
func (m *meetingRepoImpl) Append(ctx context.Context, meeting *models.Meeting) error {
	key := cache.MeetingKey(meeting.ClassID)
	log.Printf("Appending item under: %s", key)

	_, err := m.client.Set(ctx, key, meeting.ToBytes(), time.Duration(time.Minute*60)).Result()
	if err != nil {
		return fmt.Errorf("failed to append new value with key:%s", meeting.ClassID)
	}
	return nil
}

// Delete implements repository.MeetingRepository.
func (m *meetingRepoImpl) Delete(ctx context.Context, classID string) error {
	key := cache.MeetingKey(classID)
	removedAmount, err := m.client.Del(ctx, key).Result()
	if err != nil || removedAmount < 1 {
		log.Println(err)
		return fmt.Errorf("failed to delete value with id: %s", classID)
	}

	return nil
}

// Get implements repository.MeetingRepository.
func (m *meetingRepoImpl) Get(ctx context.Context, classID string) (*dto.MeetingDTO, error) {
	key := cache.MeetingKey(classID)
	log.Printf("Getting key: %s", key)

	result, err := m.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s not found in Redis", key)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	var meeting models.Meeting
	if err := json.Unmarshal([]byte(result), &meeting); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meeting for key %s: %w", key, err)
	}

	log.Printf("Found meeting for key %s", key)
	return meeting.DTO(), nil
}

func NewMeetingRepository(client *redis.Client) repository.MeetingRepository {
	return &meetingRepoImpl{
		client: client,
	}
}
