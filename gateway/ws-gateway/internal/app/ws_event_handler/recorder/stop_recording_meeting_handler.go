package recorder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/domain/broker"
	"ws-gateway/internal/domain/event"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/recorder"
)

type stopRecordMeetingHandler struct {
	broker       interfaces.Broker
	exchange     string
	hubManager   interfaces.HubManager
	cacheService interfaces.CacheEventService
}

// Handle implements interfaces.EventHandler.
func (s *stopRecordMeetingHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var evt recorder.StopRecordingRequestedWSEvent
	var wg sync.WaitGroup
	if err := json.Unmarshal(body, &evt); err != nil {
		return fmt.Errorf("failed to decode event: [%s]", evt.Name())
	}

	newEvent := &event.StopRecordingMeetingEvent{
		StopTime: evt.StopTime,
		RoomID:   evt.RoomID,
	}

	var errorsCh chan error = make(chan error, 2)

	wg.Go(func() {
		dest := broker.NewExchangeDestination(newEvent, s.exchange)

		err := s.broker.Publish(ctx, newEvent, dest)
		if err != nil {
			errorsCh <- fmt.Errorf("failed to publish %s", newEvent.Name())
		}
	})

	wg.Go(func() {
		err := s.cacheService.RemoveMeetingFromPool(ctx, evt.RoomID)
		if err != nil {
			errorsCh <- err
		}
	})

	wg.Go(func() {
		var err error
		var payload []byte
		payload, err = wsevent.EncodeSocketEventWrapper(&evt)
		if err != nil {
			errorsCh <- err
			return
		}
		s.hubManager.Emit(evt.RoomID, payload, func(id string) bool { return true })
	})

	wg.Wait()

	close(errorsCh)

	if len(errorsCh) > 0 {
		var errs []error
		for err := range errorsCh {
			errs = append(errs, err)
		}

		return fmt.Errorf("multiple errors occurred: %w", errors.Join(errs...))
	}

	return nil
}

func NewStopRecordMeetingHandler(broker interfaces.Broker, exchange string, hubManager interfaces.HubManager, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	return &stopRecordMeetingHandler{broker: broker, exchange: exchange, hubManager: hubManager, cacheService: cacheService}
}
