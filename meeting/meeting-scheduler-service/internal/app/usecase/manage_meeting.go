package usecase

import (
	"context"
	"fmt"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/event"
	"meeting-scheduler-service/internal/domain/repository"
	"meeting-scheduler-service/internal/infrastructure/config"
	"os"
)

type manageMeetingImlp struct {
	Broker            interfaces.Broker
	MeetingRepository repository.MeetingRepository

	notificationChannelName string
	meetingChannelName      string
}

// Start implements ManageMeeting.
func (m *manageMeetingImlp) Start(startedMeetingDto dto.StartMeetingDTO) (*dto.MeetingDTO, error) {
	ctx := context.Background()

	ev := event.NewMeetingStartedEvent(startedMeetingDto)

	err := m.MeetingRepository.Append(ctx, startedMeetingDto.ClassID, ev.StartedTime)
	if err != nil {
		return nil, err
	}

	err = m.Broker.Publish(ev, m.meetingChannelName, m.notificationChannelName)
	if err != nil {
		return nil, fmt.Errorf("failed to create meeting, try again")
	}

	startedTime := ev.StartedTime.Unix()

	result := dto.NewMeetingDTO(ev.MeetingID, ev.Members, &startedTime)

	return result, nil
}

// Stop implements ManageMeeting.
func (m *manageMeetingImlp) Stop(endMeetingDto dto.EndMeetingDTO) (*dto.MeetingDTO, error) {
	ctx := context.Background()

	ev := event.NewMeetingEndedEvent(endMeetingDto)
	err := m.MeetingRepository.Delete(ctx, endMeetingDto.ClassID)
	if err != nil {
		return nil, err
	}

	err = m.Broker.Publish(ev, m.meetingChannelName)

	if err != nil {
		return nil, fmt.Errorf("Failed to stop meeting, try again")
	}

	result := dto.NewMeetingDTO(endMeetingDto.MeetingID, endMeetingDto.Members, nil)

	return result, nil
}

func NewMeetingManager(broker interfaces.Broker, repo repository.MeetingRepository) interfaces.ManageMeeting {
	meetingChannelName := os.Getenv(config.EVENT_EXCHANGE_QUEUE_NAME)
	notificationChannelName := os.Getenv(config.NOTIFICATION_EXCHANGE_QUEUE_NAME)

	return &manageMeetingImlp{
		Broker:            broker,
		MeetingRepository: repo,

		meetingChannelName:      meetingChannelName,
		notificationChannelName: notificationChannelName,
	}
}
