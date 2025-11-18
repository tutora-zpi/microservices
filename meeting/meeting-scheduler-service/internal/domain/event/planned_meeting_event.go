package event

import (
	"meeting-scheduler-service/internal/domain/dto"
	"reflect"
)

type PlannedMeetingEvent struct {
	dto.PlanMeetingDTO
}

func (p *PlannedMeetingEvent) Name() string {
	return reflect.TypeOf(*p).Name()
}

func NewPlannedMeetingEvent(dto dto.PlanMeetingDTO) Event {
	return &PlannedMeetingEvent{dto}
}
