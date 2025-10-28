package cache

import (
	"context"
	"fmt"
	"log"
	"recorder-service/internal/domain/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type botRepoImpl struct {
	client *redis.Client
}

// TryAdd implements repository.BotRepository.
func (b *botRepoImpl) TryAdd(ctx context.Context, roomID string, botID string, ttl time.Duration) error {
	key := BotKey(roomID)
	hasAdded, err := b.client.SetNX(ctx, key, botID, ttl).Result()
	if err != nil {
		log.Printf("Failed to add bot to cache: %v", err)
		return fmt.Errorf("failed to add bot %s to room: %s", botID, roomID)
	}

	if !hasAdded {
		return fmt.Errorf("room is already added")
	}

	return nil
}

// Delete implements repository.BotRepository.
func (b *botRepoImpl) Delete(ctx context.Context, roomID string) error {
	key := BotKey(roomID)
	_, err := b.client.Del(ctx, key).Result()
	if err != nil {
		log.Printf("Failed to delete room from cache: %v", err)
		return fmt.Errorf("failed to delete room from cache %s", roomID)
	}

	return nil
}

func NewBotRepository(client *redis.Client) repository.BotRepository {
	return &botRepoImpl{client: client}
}
