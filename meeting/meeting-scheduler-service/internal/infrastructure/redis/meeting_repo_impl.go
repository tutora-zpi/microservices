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
		return fmt.Errorf("failed to delete value with id:%s", classID)
	}

	return nil
}

// Get implements repository.MeetingRepository.
func (m *meetingRepoImpl) Get(ctx context.Context, classID string) (*dto.MeetingDTO, error) {
	key := m.key(classID)
	log.Printf("Getting key: %s\n", key)

	result, err := m.client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key %s not found in Redis\n", key)
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	var meeting models.Meeting
	if err := json.Unmarshal([]byte(result), &meeting); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meeting for key %s: %w", key, err)
	}

	log.Printf("Found meeting for key %s\n", key)
	return meeting.ToDTO(), nil
}

func NewMeetingRepo() repository.MeetingRepository {
	opts := &redis.Options{
		Addr:     os.Getenv(config.REDIS_ADDR),
		Password: os.Getenv(config.REDIS_PASSWORD),
		DB:       0,
	}

	client := redis.NewClient(opts)

	pong, err := client.Ping(context.Background()).Result()
	if err != nil || pong != "PONG" {
		log.Panicf("Failed to ping Redis: %v", err)
	}
	log.Println("Successfully connected to Redis:", pong)

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
