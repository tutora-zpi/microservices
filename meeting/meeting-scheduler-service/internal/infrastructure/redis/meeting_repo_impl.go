package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/models"
	"meeting-scheduler-service/internal/domain/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type meetingRepoImpl struct {
	client           *redis.Client
	temporaryMeeting func(suffix string) string
}

// Exists implements repository.MeetingRepository.
func (m *meetingRepoImpl) Exists(ctx context.Context, classID string) bool {
	key := m.temporaryMeeting(classID)

	foundNumber, _ := m.client.Exists(ctx, key).Result()
	return foundNumber > 0
}

// Append implements repository.MeetingRepository.
func (m *meetingRepoImpl) Append(ctx context.Context, meeting *models.Meeting) error {
	key := m.temporaryMeeting(meeting.ClassID)
	log.Printf("Appending item under: %s", key)

	_, err := m.client.Set(ctx, key, meeting.Json(), time.Duration(time.Minute*60)).Result()
	if err != nil {
		return fmt.Errorf("failed to append new value with key:%s", meeting.ClassID)
	}
	return nil
}

// Delete implements repository.MeetingRepository.
func (m *meetingRepoImpl) Delete(ctx context.Context, classID string) error {
	key := m.temporaryMeeting(classID)
	removedAmount, err := m.client.Del(ctx, key).Result()
	if err != nil || removedAmount < 1 {
		return fmt.Errorf("failed to delete value with id:%s", classID)
	}

	return nil
}

// Get implements repository.MeetingRepository.
func (m *meetingRepoImpl) Get(ctx context.Context, classID string) (*dto.MeetingDTO, error) {
	key := m.temporaryMeeting(classID)
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

func NewMeetingRepo(ctx context.Context, redisConfig RedisConfig) (repository.MeetingRepository, error) {
	opts := &redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	}

	client := redis.NewClient(opts)

	pong, err := client.Ping(ctx).Result()
	if err != nil || pong != "PONG" {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	log.Println("Successfully connected to Redis:", pong)

	return &meetingRepoImpl{
		client: client,
		temporaryMeeting: func(suffix string) string {
			return fmt.Sprintf("meeting:%s", suffix)
		},
	}, nil
}

func (m *meetingRepoImpl) Close() {
	if err := m.client.Close(); err != nil {
		log.Printf("Failed to close redis client: %v", err)
		return
	}

	log.Println("Redis Successfully closed connection.")
}
