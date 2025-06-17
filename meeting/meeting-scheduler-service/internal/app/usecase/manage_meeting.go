package usecase

import (
	"fmt"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/event"
)

type manageMeetingImlp struct {
	Broker interfaces.Broker
	// instace redis
}

// Start implements ManageMeeting.
func (m *manageMeetingImlp) Start(dto dto.StartMeetingDTO) (*dto.MeetingDTO, error) {
	ev := event.NewMeetingStartedEvent(dto)

	err := m.Broker.Publish(*ev)
	if err != nil {
		return nil, fmt.Errorf("Failed to create meeting, try again")
	}

	res := event.NewMeetingDTO(*ev)

	return &res, nil
}

// Stop implements ManageMeeting.
func (m *manageMeetingImlp) Stop(dto dto.EndMeetingDTO) (*dto.MeetingDTO, error) {
	ev := event.NewMeetingEndedEvent(dto)
	err := m.Broker.Publish(*ev)
	if err != nil {
		return nil, fmt.Errorf("Failed to stop meeting, try again")
	}

	res := event.NewMeetingDTO(*ev)

	return &res, nil
}

func NewMeetingManager(b interfaces.Broker) interfaces.ManageMeeting {
	return &manageMeetingImlp{
		Broker: b,
	}
}
