package dto

import (
	"time"
)

// MeetingDTO represents meeting details returned in responses.
// swagger:model MeetingDTO
type MeetingDTO struct {
	// Meeting unique identifier
	MeetingID string `json:"meetingId"`

	// Members who participated in the meeting
	Members []UserDTO `json:"members,omitempty"`

	// Meeting's started time
	StartedDate *time.Time `json:"startedTime,omitempty"`

	// Meeting's finish time
	FinishDate time.Time `json:"finishDate" validate:"required" example:"2025-10-10T12:50:05+02:00"`

	// Meetings title
	Title string `json:"title"`
}

func NewMeetingDTO(meetingID string, members []UserDTO, startedDate, finishDate time.Time, title string) *MeetingDTO {
	return &MeetingDTO{
		MeetingID:   meetingID,
		Members:     members,
		StartedDate: &startedDate,
		Title:       title,
		FinishDate:  finishDate,
	}
}

func (dto *MeetingDTO) GetStartMeetingDTO(classID string) *StartMeetingDTO {
	return &StartMeetingDTO{
		Members: dto.Members,
		ClassID: classID,
		Title:   dto.Title,
	}
}
