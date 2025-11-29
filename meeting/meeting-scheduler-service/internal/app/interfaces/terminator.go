package interfaces

import (
	"context"
	"meeting-scheduler-service/internal/domain/dto"
)

type MeetingTerminator interface {
	AppendNewMeeting(endMeetingDTO dto.EndMeetingDTO, expectedEndTimestamp int64) error
	Run(ctx context.Context, stopHandler func(ctx context.Context, endMeetingDTO dto.EndMeetingDTO) error)
}
