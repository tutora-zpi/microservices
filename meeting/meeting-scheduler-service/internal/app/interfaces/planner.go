package interfaces

import (
	"context"
)

type MeetingPlanner interface {
	LoadScheduledMeetings(ctx context.Context) error
	Listen(ctx context.Context)
	RerunNotStartedMeetings(ctx context.Context) error
}
