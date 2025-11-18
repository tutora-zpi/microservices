package event

import (
	"reflect"
	"ws-gateway/internal/domain/ws_event/board"
)

type BoardUpdateEvent struct {
	board.BoardUpdateWSEvent
}

func NewBoardUpdateEvent(wsevent board.BoardUpdateWSEvent) *BoardUpdateEvent {
	return &BoardUpdateEvent{
		BoardUpdateWSEvent: wsevent,
	}
}

func (b *BoardUpdateEvent) Name() string {
	return reflect.TypeOf(*b).Name()
}
