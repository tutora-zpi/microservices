package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/models"
	"meeting-scheduler-service/internal/domain/repository"
	"meeting-scheduler-service/internal/infrastructure/config"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type meetingRepoImpl struct {
	client *redis.Client
	key    func(suffix string) string
}

// Append implements repository.MeetingRepository.
func (m *meetingRepoImpl) Append(ctx context.Context, meeting *models.Meeting) error {
	key := m.key(meeting.ClassID)
	log.Printf("Appending item under: %s\n", key)

	_, err := m.client.Set(ctx, key, meeting.ToJSON(), time.Duration(time.Minute*70)).Result()
	if err != nil {
		return fmt.Errorf("failed to append new value with key:%s", meeting.ClassID)
	}
	return nil
}

// Delete implements repository.MeetingRepository.
func (m *meetingRepoImpl) Delete(ctx context.Context, classID string) error {
	key := m.key(classID)
	removedAmount, err := m.client.Del(ctx, key).Result()
	if err != nil || removedAmount < 1 {
		return fmt.Errorf("failed to delete value with key:%s", classID)
	}

	return nil
}

// Get implements repository.MeetingRepository.
func (m *meetingRepoImpl) Get(ctx context.Context, classID string) (*dto.MeetingDTO, error) {
	key := m.key(classID)
	result, err := m.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var meeting *models.Meeting

	if err := json.Unmarshal([]byte(result), &meeting); err != nil {
		return nil, err
	}

	return meeting.ToDTO(), nil
}

func NewMeetingRepo() repository.MeetingRepository {
	ops := &redis.Options{
		Addr:     os.Getenv(config.REDIS_ADDR),
		Password: os.Getenv(config.REDIS_PASSWORD),
		DB:       0,
	}

	client := redis.NewClient(ops)

	result, err := client.Ping(context.Background()).Result()
	if err != nil || result != "PONG" {
		log.Panicf("Failed to connect with Redis")
	} else {
		log.Println("Successfully connected to Redis")
	}

	return &meetingRepoImpl{
		client: client,
		key: func(suffix string) string {
			return fmt.Sprintf("meeting:%s", suffix)
		},
	}
}

func (m *meetingRepoImpl) Close() {
	if err := m.client.Close(); err != nil {
		log.Printf("Failed to close redis client: %v\n", err)
		return
	}

	log.Println("Successfully closed connection.")
}
