package usecase

import (
	"context"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/broker"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/event"
	"meeting-scheduler-service/internal/domain/repository"
	"time"
)

type ManageMeetingImlp struct {
	Broker                    interfaces.Broker
	MeetingRepository         repository.MeetingRepository
	PlannedMeetingsRepository repository.PlannedMeetingsRepository

	NotificationExchange string
	MeetingExchange      string
}

// CancelPlannedMeeting implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) CancelPlannedMeeting(ctx context.Context, id int) error {
	panic("unimplemented")
}

// GetPlannedMeetings implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) GetPlannedMeetings(ctx context.Context, dto dto.FetchPlannedMeetings) ([]dto.PlannedMeetingDTO, error) {
	results, err := m.PlannedMeetingsRepository.GetPlannedMeetings(ctx, dto)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// LoadMorePlannedMeetings implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) LoadMorePlannedMeetings(ctx context.Context, interval time.Duration) ([]dto.PlanMeetingDTO, error) {
	// time window
	start := time.Now().UTC().Add(-time.Minute * 2)
	before := time.Now().UTC().Add(interval)

	results, err := m.PlannedMeetingsRepository.ProcessPlannedMeetings(ctx, start, before)

	return results, err
}

// Plan implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) Plan(ctx context.Context, dto dto.PlanMeetingDTO) (*dto.PlanMeetingDTO, error) {
	log.Printf("Planning meeting on %s", dto.StartDate.Format(time.RFC3339))

	if m.MeetingRepository.Exists(ctx, dto.ClassID) {
		return nil, fmt.Errorf("meeting has already started")
	}

	if !m.PlannedMeetingsRepository.CanStartAnotherMeeting(ctx, dto) {
		return nil, fmt.Errorf("unable to start another cause of meeting collisions")
	}

	createdMeeting, err := m.PlannedMeetingsRepository.CreatePlannedMeetings(ctx, dto)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create planned meeting")
	}

	log.Printf("Meeting with id: %s successfully planned\n", dto.ClassID)

	plannedMeetingEvent := event.NewPlannedMeetingEvent(dto)
	dest := broker.NewExchangeDestination(plannedMeetingEvent, m.NotificationExchange)
	if err := m.Broker.Publish(ctx, plannedMeetingEvent, dest); err != nil {
		log.Printf("Failed to publish notification: %v\n", err)
	}

	return createdMeeting, nil
}

// ActiveMeeting implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) ActiveMeeting(ctx context.Context, classID string) (*dto.MeetingDTO, error) {
	return m.MeetingRepository.Get(ctx, classID)
}

// Start implements ManageMeeting.
func (m *ManageMeetingImlp) Start(ctx context.Context, startedMeetingDto dto.StartMeetingDTO) (*dto.MeetingDTO, error) {
	log.Println("Starting meeting")
	if m.MeetingRepository.Exists(ctx, startedMeetingDto.ClassID) {
		return nil, fmt.Errorf("meeting has already started")
	}

	meetingStartedEvent := event.NewMeetingStartedEvent(startedMeetingDto)

	meeting := meetingStartedEvent.NewMeeting(startedMeetingDto.ClassID, startedMeetingDto.Title)

	err := m.MeetingRepository.Append(ctx, meeting)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to add new meeting to cache")
	}

	log.Printf("Appended to cache new meeting in class: %s\n", meeting.ClassID)

	destinations := broker.NewMultipleDestination(meetingStartedEvent, m.MeetingExchange, m.NotificationExchange)

	err = m.Broker.PublishMultiple(
		ctx,
		meetingStartedEvent,
		destinations...,
	)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create meeting, try again")
	}

	log.Println("Successfully published events")

	result := dto.NewMeetingDTO(meetingStartedEvent.MeetingID, meetingStartedEvent.Members, &meeting.Timestamp, meeting.Title)

	return result, nil
}

// Stop implements ManageMeeting.
func (m *ManageMeetingImlp) Stop(ctx context.Context, endMeetingDto dto.EndMeetingDTO) error {
	meetingEndedEvent := event.NewMeetingEndedEvent(endMeetingDto)
	err := m.MeetingRepository.Delete(ctx, endMeetingDto.ClassID)
	if err != nil {
		return err
	}

	err = m.Broker.Publish(ctx, meetingEndedEvent, broker.NewExchangeDestination(meetingEndedEvent, m.MeetingExchange))

	if err != nil {
		return fmt.Errorf("failed to stop meeting, try again")
	}

	log.Println("Meeting has been finished")

	return nil
}

func NewManageMeeting(
	broker interfaces.Broker,
	meetingRepo repository.MeetingRepository,
	plannedMeetingsRepo repository.PlannedMeetingsRepository,
	notificationExchange string,
	meetingExchange string,
) interfaces.ManageMeeting {
	return &ManageMeetingImlp{
		Broker:                    broker,
		MeetingRepository:         meetingRepo,
		PlannedMeetingsRepository: plannedMeetingsRepo,
		NotificationExchange:      notificationExchange,
		MeetingExchange:           meetingExchange,
	}
}
