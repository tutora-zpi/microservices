package repository

import (
	"context"
	"time"
)

type MeetingRepository interface {
	Append(ctx context.Context, classID string, timestamp time.Time) error
	Contains(ctx context.Context, classID string) bool
	Delete(ctx context.Context, classID string) error
	Close()
}
