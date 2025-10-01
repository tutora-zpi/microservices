package redis

import (
	"context"
	"fmt"
	"log"
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
func (m *meetingRepoImpl) Append(ctx context.Context, classID string, timestamp time.Time) error {
	key := m.key(classID)
	_, err := m.client.Set(ctx, key, timestamp.UnixNano(), time.Duration(time.Minute*70)).Result()
	if err != nil {
		return fmt.Errorf("failed to append new value with key:%s", classID)
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
func (m *meetingRepoImpl) Contains(ctx context.Context, classID string) bool {
	key := m.key(classID)
	_, err := m.client.Get(ctx, key).Result()
	return err == nil
}

func NewMeetingRepo() repository.MeetingRepository {
	ops := &redis.Options{
		Addr:     os.Getenv(config.REDIS_ADDR),
		Password: os.Getenv(config.REDIS_PASSWORD),
		DB:       0,
	}

	client := redis.NewClient(ops)

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
