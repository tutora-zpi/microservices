package interfaces

import "meeting-scheduler-service/internal/domain/dto"

type ManageMeeting interface {
	Start(dto dto.StartMeetingDTO) (*dto.MeetingDTO, error)
	Stop(dto dto.EndMeetingDTO) (*dto.MeetingDTO, error)
}
