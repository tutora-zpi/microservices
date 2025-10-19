package repo

import (
	"context"
	"fmt"
	"time"
	"ws-gateway/internal/infrastructure/cache/enum"

	"github.com/redis/go-redis/v9"
)

type CacheEventRepository interface {
	PushEvent(ctx context.Context, roomID string, compressedData []byte) error
	SaveSnapshot(ctx context.Context, roomID string, compressedData []byte) error
	GetSnapshot(ctx context.Context, roomID string) ([]byte, error)
	GetCachedEvents(ctx context.Context, roomID string) ([][]byte, error)
}

type cacheEventRepoImpl struct {
	client     *redis.Client
	maxPerRoom int
	ttl        time.Duration
}

// GetSnapshot implements CacheEventRepository.
func (c *cacheEventRepoImpl) GetSnapshot(ctx context.Context, roomID string) ([]byte, error) {
	key := enum.SnapshotKey(roomID)

	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("something went wrong: %v", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("not found: len(result) is 0")
	}

	return []byte(result), nil
}

// SaveSnapshot implements CacheEventRepository.
func (c *cacheEventRepoImpl) SaveSnapshot(ctx context.Context, roomID string, compressedData []byte) error {
	key := enum.SnapshotKey(roomID)
	_, err := c.client.Set(ctx, key, compressedData, c.ttl).Result()
	if err != nil {
		return fmt.Errorf("failed to save event: %v", err)
	}
	return nil
}

// PushEvent implements interfaces.CacheEventRepository.
func (c *cacheEventRepoImpl) PushEvent(ctx context.Context, roomID string, compressedData []byte) error {
	key := enum.EventKey(roomID)

	pipe := c.client.Pipeline()

	pipe.LPush(ctx, key, compressedData)
	pipe.LTrim(ctx, key, 0, int64(c.maxPerRoom)-1)
	pipe.Expire(ctx, key, c.ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to push event with trim: %v", err)
	}
	return nil
}

// GetCachedEvents implements interfaces.CacheEventRepository.
func (c *cacheEventRepoImpl) GetCachedEvents(ctx context.Context, roomID string) ([][]byte, error) {
	key := enum.EventKey(roomID)
	results, err := c.client.LRange(ctx, key, 0, -1).Result()

	if err != nil {
		return nil, fmt.Errorf("something went wrong: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no cached events for key: %s", roomID)
	}

	byteMatrix := make([][]byte, len(results))
	for i, r := range results {
		byteMatrix[i] = []byte(r)
	}

	return byteMatrix, nil
}

func NewCacheEventRepository(client *redis.Client, maxPerRoom int, ttl time.Duration) CacheEventRepository {
	return &cacheEventRepoImpl{
		client:     client,
		maxPerRoom: maxPerRoom,
		ttl:        ttl,
	}
}
