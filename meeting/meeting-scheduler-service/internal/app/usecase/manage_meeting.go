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
	broker                    interfaces.Broker
	meetingRepository         repository.MeetingRepository
	plannedMeetingsRepository repository.PlannedMeetingsRepository

	meetingExchange string
}

// CancelPlannedMeeting implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) CancelPlannedMeeting(ctx context.Context, id int) error {
	log.Printf("Cancelling planned meeting with: %d", id)
	if err := m.plannedMeetingsRepository.CancelMeeting(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetPlannedMeetings implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) GetPlannedMeetings(ctx context.Context, dto dto.FetchPlannedMeetingsDTO) ([]dto.PlannedMeetingDTO, error) {
	log.Printf("Getting planned meetings [%v]\n", dto)

	results, err := m.plannedMeetingsRepository.GetPlannedMeetings(ctx, dto)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// LoadMorePlannedMeetings implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) LoadMorePlannedMeetings(ctx context.Context, interval time.Duration) ([]dto.PlanMeetingDTO, error) {
	// time window
	now := time.Now().UTC()
	start := now.Add(-time.Minute)
	before := now.Add(interval)

	results, err := m.plannedMeetingsRepository.ProcessPlannedMeetings(ctx, start, before)

	return results, err
}

// Plan implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) Plan(ctx context.Context, dto dto.PlanMeetingDTO) (*dto.PlannedMeetingDTO, error) {
	log.Printf("Planning meeting on %s", dto.StartDate.Format(time.RFC3339))

	if m.meetingRepository.Exists(ctx, dto.ClassID) {
		return nil, fmt.Errorf("meeting has already started")
	}

	if !m.plannedMeetingsRepository.CanStartAnotherMeeting(ctx, dto) {
		return nil, fmt.Errorf("unable to start another cause of meeting collisions")
	}

	createdMeeting, err := m.plannedMeetingsRepository.CreatePlannedMeetings(ctx, dto)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create planned meeting")
	}

	log.Printf("Meeting with id: %s successfully planned\n", dto.ClassID)

	plannedMeetingEvent := event.NewPlannedMeetingEvent(dto)
	dest := broker.NewExchangeDestination(plannedMeetingEvent, m.meetingExchange)
	if err := m.broker.Publish(ctx, plannedMeetingEvent, dest); err != nil {
		log.Printf("Failed to publish notification: %v\n", err)
	}

	return createdMeeting, nil
}

// ActiveMeeting implements interfaces.ManageMeeting.
func (m *ManageMeetingImlp) ActiveMeeting(ctx context.Context, classID string) (*dto.MeetingDTO, error) {
	log.Printf("Checking for active meetings in class: %s\n", classID)
	return m.meetingRepository.Get(ctx, classID)
}

// Start implements ManageMeeting.
func (m *ManageMeetingImlp) Start(ctx context.Context, startedMeetingDto dto.StartMeetingDTO) (*dto.MeetingDTO, error) {
	log.Println("Starting meeting")
	if m.meetingRepository.Exists(ctx, startedMeetingDto.ClassID) {
		return nil, fmt.Errorf("meeting has already started")
	}

	meetingStartedEvent := event.NewMeetingStartedEvent(startedMeetingDto)

	meeting := meetingStartedEvent.NewMeeting(startedMeetingDto)

	err := m.meetingRepository.Append(ctx, meeting)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to add new meeting to cache")
	}

	log.Printf("Appended to cache new meeting in class: %s\n", meeting.ClassID)

	destination := broker.NewExchangeDestination(meetingStartedEvent, m.meetingExchange)

	err = m.broker.Publish(
		ctx,
		meetingStartedEvent,
		destination,
	)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create meeting, try again")
	}

	log.Println("Successfully published events")

	result := dto.NewMeetingDTO(
		meetingStartedEvent.MeetingID,
		meetingStartedEvent.Members,
		meetingStartedEvent.StartedTime,
		meetingStartedEvent.FinishTime,
		meeting.Title,
	)

	return result, nil
}

// Stop implements ManageMeeting.
func (m *ManageMeetingImlp) Stop(ctx context.Context, endMeetingDto dto.EndMeetingDTO) error {
	log.Println("Stopping meeting...")

	meetingEndedEvent := event.NewMeetingEndedEvent(endMeetingDto)
	err := m.meetingRepository.Delete(ctx, endMeetingDto.ClassID)
	if err != nil {
		return err
	}

	err = m.broker.Publish(ctx, meetingEndedEvent, broker.NewExchangeDestination(meetingEndedEvent, m.meetingExchange))

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
	meetingExchange string,
) interfaces.ManageMeeting {
	return &ManageMeetingImlp{
		broker:                    broker,
		meetingRepository:         meetingRepo,
		plannedMeetingsRepository: plannedMeetingsRepo,
		meetingExchange:           meetingExchange,
	}
}
