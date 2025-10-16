package repository

import (
	"context"
	"fmt"
	"log"
	"signaling-service/internal/domain/enum"
	"signaling-service/internal/domain/models"
	"signaling-service/internal/domain/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type statusRepoImpl struct {
	client          *redis.Client
	temporaryStatus func(suffix string) string
}

// Delete implements repository.StatusRepository.
func (s *statusRepoImpl) Delete(ctx context.Context, userID string) error {
	key := s.temporaryStatus(userID)
	removedAmount, err := s.client.Del(ctx, key).Result()
	if err != nil || removedAmount < 1 {
		return fmt.Errorf("failed to delete value with id:%s", userID)
	}

	return nil
}

// Get implements repository.StatusRepository.
func (s *statusRepoImpl) Get(ctx context.Context, userID string) (*models.Status, error) {
	key := s.temporaryStatus(userID)
	log.Printf("Getting key: %s", key)

	result, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s not found in Redis", key)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	model, err := models.DecodeStatus([]byte(result))
	if err != nil {
		return nil, err
	}

	log.Printf("Found status for key %s", key)
	return model, nil
}

// Save implements repository.StatusRepository.
func (s *statusRepoImpl) Save(ctx context.Context, userID string, status enum.UserStatus, ttl time.Duration) error {
	key := s.temporaryStatus(userID)

	body, err := models.EncodeStatus(userID, status)
	if err != nil {
		return err
	}

	_, err = s.client.Set(ctx, key, body, ttl).Result()

	if err != nil {
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return nil
}

func NewStatusRepository(client *redis.Client) repository.StatusRepository {
	return &statusRepoImpl{
		client: client,
		temporaryStatus: func(suffix string) string {
			return fmt.Sprintf("status:%s", suffix)
		},
	}
}
