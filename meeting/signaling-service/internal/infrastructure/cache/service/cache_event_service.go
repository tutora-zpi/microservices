package service

import (
	"context"
	"fmt"
	"log"
	"signaling-service/internal/app/interfaces"
	wsevent "signaling-service/internal/domain/ws_event"
	cache "signaling-service/internal/infrastructure/cache/repo"
	"signaling-service/pkg/gzip"
)

type cacheEventServiceImpl struct {
	repo cache.CacheEventRepository
}

// GetSnapshot implements interfaces.CacheEventService.
func (c *cacheEventServiceImpl) GetSnapshot(ctx context.Context, roomID string) ([]byte, error) {
	compressedData, err := c.repo.GetSnapshot(ctx, roomID)
	if err != nil {
		return nil, err
	}

	decompressed, err := gzip.Decompress(compressedData)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

// MakeSnapshot implements interfaces.CacheEventService.
func (c *cacheEventServiceImpl) MakeSnapshot(ctx context.Context, event wsevent.SocketEventWrapper, roomID string) error {
	compressedData, err := gzip.Compress(event.ToBytes())
	if err != nil {
		log.Printf("Error during compressing data: %v", err)
		return err
	}

	err = c.repo.SaveSnapshot(ctx, roomID, compressedData)

	if err != nil {
		log.Printf("Error during saving snapshot: %v", err)
		return err
	}

	return nil
}

// GetLastEventsData implements interfaces.CacheEventService.
func (c *cacheEventServiceImpl) GetLastEventsData(ctx context.Context, roomID string) ([][]byte, error) {
	log.Printf("Getting last events for %s", roomID)

	compressedStack, err := c.repo.GetCachedEvents(ctx, roomID)
	if err != nil {
		return nil, err
	}

	rowsNumber := len(compressedStack)
	decompressedData := make([][]byte, rowsNumber)

	for i, row := range compressedStack {
		decompressed, err := gzip.Decompress(row)
		if err != nil {
			log.Printf("Error during deompressing: %v", err)
		} else {
			decompressedData[i] = decompressed
		}
	}

	return decompressedData, nil
}

// PushRecentEvent implements interfaces.CacheEventService.
func (c *cacheEventServiceImpl) PushRecentEvent(ctx context.Context, event wsevent.SocketEventWrapper, roomID string) error {
	log.Printf("Pushing new event %v to room with id: %s", event, roomID)

	compressed, err := gzip.Compress(event.ToBytes())
	if err != nil {
		return fmt.Errorf("failed to compress data: %v", err)
	}

	if err := c.repo.PushEvent(ctx, roomID, compressed); err != nil {
		return err
	}

	return nil
}

func NewCacheEventSerivce(repo cache.CacheEventRepository) interfaces.CacheEventService {
	return &cacheEventServiceImpl{repo: repo}
}
