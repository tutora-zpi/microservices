package repository

import (
	"context"
	"meeting-scheduler-service/internal/domain/dto"
	"time"
)

type PlannedMeetingsRepository interface {
	CanStartAnotherMeeting(ctx context.Context, dto dto.PlanMeetingDTO) bool
	ProcessPlannedMeetings(ctx context.Context, start time.Time, before time.Time) ([]dto.PlanMeetingDTO, error)
	CreatePlannedMeetings(ctx context.Context, dto dto.PlanMeetingDTO) (*dto.PlanMeetingDTO, error)
	GetPlannedMeetings(ctx context.Context, dto dto.FetchPlannedMeetings) ([]dto.PlannedMeetingDTO, error)
	CancelMeeting(ctx context.Context, id int) error
	Close()
}
