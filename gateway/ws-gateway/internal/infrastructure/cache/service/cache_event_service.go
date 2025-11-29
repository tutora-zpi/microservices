package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"ws-gateway/internal/app/interfaces"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/recorder"
	"ws-gateway/internal/infrastructure/cache/enum"
	cache "ws-gateway/internal/infrastructure/cache/repo"
	"ws-gateway/pkg/gzip"
)

type cacheEventServiceImpl struct {
	repo cache.CacheEventRepository
}

// DeleteSnapshot implements interfaces.CacheEventService.
func (c *cacheEventServiceImpl) DeleteSnapshot(ctx context.Context, roomID string) error {
	return c.repo.Del(ctx, enum.SnapshotKey(roomID))
}

// IsMeetingRecorded implements interfaces.CacheEventService.
func (c *cacheEventServiceImpl) IsMeetingRecorded(ctx context.Context, keyRoomID string) (*recorder.RecordRequestedWSEvent, error) {
	res, err := c.repo.Get(ctx, keyRoomID)
	if err != nil {
		return nil, err
	}

	evt := recorder.NewRecordRequestedWSEventFromBytes(res)

	if evt == nil {
		return nil, fmt.Errorf("failed to decode event")
	}

	return evt, nil
}

// RemoveMeetingFromPool implements interfaces.CacheEventService.
func (c *cacheEventServiceImpl) RemoveMeetingFromPool(ctx context.Context, keyRoomID string) error {
	return c.repo.Del(ctx, enum.IsRecorded(keyRoomID))
}

// SetMeetingIsRecorded implements interfaces.CacheEventService.
func (c *cacheEventServiceImpl) SetMeetingIsRecorded(ctx context.Context, keyRoomID string, evt recorder.RecordRequestedWSEvent) error {
	ttl := time.Until(evt.FinishTime)

	return c.repo.Set(ctx, keyRoomID, ttl, evt.ToBytes())
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
	log.Printf("Getting last events for room %s", roomID)

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
	log.Printf("Pushing new event %s to room with id: %s", event.Name, roomID)

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
