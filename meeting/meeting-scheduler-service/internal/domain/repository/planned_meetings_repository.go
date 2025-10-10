package repository

import (
	"context"
	"meeting-scheduler-service/internal/domain/dto"
	"time"
)

type PlannedMeetingsRepository interface {
	CanStartAnotherMeeting(ctx context.Context, meeting dto.PlanMeetingDTO) bool
	ProcessPlannedMeetings(ctx context.Context, start time.Time, before time.Time) ([]dto.PlanMeetingDTO, error)
	CreatePlannedMeetings(ctx context.Context, meeting dto.PlanMeetingDTO) (*dto.PlanMeetingDTO, error)
	Close()
}
