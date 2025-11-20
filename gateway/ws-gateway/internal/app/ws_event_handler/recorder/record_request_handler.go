package recorder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/domain/broker"
	"ws-gateway/internal/domain/event"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/recorder"
)

type recordRequestHandler struct {
	broker       interfaces.Broker
	exchange     string
	hubManager   interfaces.HubManager
	cacheService interfaces.CacheEventService
}

// Handle implements interfaces.EventHandler.
func (r *recordRequestHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var evt recorder.RecordRequestedWSEvent
	var wg sync.WaitGroup
	if err := json.Unmarshal(body, &evt); err != nil {
		return fmt.Errorf("failed to decode event: [%s]", evt.Name())
	}

	expectedUsers := r.hubManager.GetUsersFromRoomID(evt.RoomID)
	if len(expectedUsers) < 1 {
		return fmt.Errorf("no people in the room: %s", evt.RoomID)
	}
	log.Printf("Expected user to be recorded: %v", expectedUsers)

	newEvent := event.NewRecordMeetingEvent(evt, expectedUsers)

	var errorsCh chan error = make(chan error, 3)

	wg.Go(func() {
		dest := broker.NewExchangeDestination(newEvent, r.exchange)

		err := r.broker.Publish(ctx, newEvent, dest)
		if err != nil {
			errorsCh <- fmt.Errorf("failed to publish %s", newEvent.Name())
		}
	})

	wg.Go(func() {
		err := r.cacheService.SetMeetingIsRecorded(ctx, evt.RoomID, evt)
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
		r.hubManager.Emit(evt.RoomID, payload, func(id string) bool { return true })
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

func NewRecordRequestHandler(broker interfaces.Broker, exchange string, hubManager interfaces.HubManager, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	return &recordRequestHandler{broker: broker, exchange: exchange, hubManager: hubManager, cacheService: cacheService}
}
