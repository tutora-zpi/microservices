package interfaces

import (
	"context"
	"meeting-scheduler-service/internal/domain/dto"
	"time"
)

type ManageMeeting interface {
	Start(ctx context.Context, dto dto.StartMeetingDTO) (*dto.MeetingDTO, error)
	Stop(ctx context.Context, dto dto.EndMeetingDTO) error
	ActiveMeeting(ctx context.Context, classID string) (*dto.MeetingDTO, error)
	Plan(ctx context.Context, dto dto.PlanMeetingDTO) (*dto.PlannedMeetingDTO, error)
	LoadMorePlannedMeetings(ctx context.Context, interval time.Duration) ([]dto.PlanMeetingDTO, error)
	GetPlannedMeetings(ctx context.Context, dto dto.FetchPlannedMeetingsDTO) ([]dto.PlannedMeetingDTO, error)
	CancelPlannedMeeting(ctx context.Context, id int) error
}
