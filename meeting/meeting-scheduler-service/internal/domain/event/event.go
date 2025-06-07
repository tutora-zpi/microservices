package event

import (
	"encoding/json"
	"log"
	"meeting-scheduler-service/internal/domain/dto"
	"reflect"
)

type EventWrapper struct {
	Pattern string `json:"pattern"`
	Data    event  `json:"data"`
}

type event any

func NewEventWrapper(pattern string, event event) *EventWrapper {
	return &EventWrapper{
		Pattern: pattern,
		Data:    event,
	}
}

func (e *EventWrapper) ToJson() ([]byte, error) {
	return json.Marshal(e)
}

func NewMeetingDTO(event EventWrapper) dto.MeetingDTO {
	log.Printf("Pattern %s\n", event.Pattern)

	switch event.Pattern {
	case reflect.TypeOf(MeetingStartedEvent{}).Name():
		if start, ok := event.Data.(MeetingStartedEvent); ok {
			return dto.MeetingDTO{
				MeetingID: start.MeetingID,
				Members:   start.Members,
			}
		}
	case reflect.TypeOf(MeetingEndedEvent{}).Name():
		log.Println("Mapping from MeetingEndedEvent to MeetingDTO")
		if end, ok := event.Data.(MeetingEndedEvent); ok {
			log.Println("Successfully mapped!")
			return dto.MeetingDTO{
				MeetingID: end.MeetingID,
				Members:   end.Members,
			}
		}
	}
	return dto.MeetingDTO{}
}
