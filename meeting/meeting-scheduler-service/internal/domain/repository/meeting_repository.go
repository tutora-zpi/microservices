package repository

import (
	"context"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/models"
)

type MeetingRepository interface {
	Append(ctx context.Context, meeting *models.Meeting) error
	Get(ctx context.Context, classID string) (*dto.MeetingDTO, error)
	Delete(ctx context.Context, classID string) error
	Exists(ctx context.Context, classID string) bool
}
