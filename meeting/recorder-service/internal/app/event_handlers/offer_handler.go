package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"recorder-service/internal/app/interfaces/handler"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/ws_event/rtc"
	"recorder-service/internal/infrastructure/webrtc/writer"
)

type offerHandler struct {
	botService    service.BotService
	writerFactory writer.WriterFactory
}

// Handle implements handler.EventHandler.
func (o *offerHandler) Handle(ctx context.Context, body []byte) error {
	log.Println("OFFER HANDLER")
	var evt rtc.OfferWSEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode Offer event: %v", err)
		return fmt.Errorf("failed to decode event: %w", err)
	}

	log.Printf("Received event: FROM=%s TO=%s ROOMID=%s", evt.From, evt.To, evt.RoomID)

	// if err := evt.IsValid(); err != nil {
	// 	return err
	// }

	bot, ok := o.botService.GetBot(evt.RoomID)
	if !ok {
		return fmt.Errorf("bot %s not found in room %s", bot.ID(), evt.RoomID)
	}

	if evt.To != bot.ID() {
		return fmt.Errorf("offer not for bot %s, offer is for %s", bot.Name(), evt.To)
	}

	peer, _ := bot.GetPeer(evt.From)

	if err := peer.SetRemoteDescription(evt.Offer); err != nil {
		log.Printf("Failed to set remote description for user %s: %v", evt.From, err)
		return err
	}

	if err := peer.CreateAnswer(); err != nil {
		log.Printf("Failed to create/send answer for user %s: %v", evt.From, err)
		return err
	}

	log.Printf("Offer handled for user %s in room %s", evt.From, evt.RoomID)
	return nil
}

func NewOfferHandler(botService service.BotService, writerFactory writer.WriterFactory) handler.EventHandler {
	return &offerHandler{
		botService:    botService,
		writerFactory: writerFactory,
	}
}
