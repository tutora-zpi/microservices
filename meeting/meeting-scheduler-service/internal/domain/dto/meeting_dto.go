package dto

// MeetingDTO represents meeting details returned in responses.
// swagger:model MeetingDTO
type MeetingDTO struct {
	// Meeting unique identifier
	MeetingID string `json:"meetingID"`
	// Members who participated in the meeting
	Members []UserDTO `json:"members"`
	// Timestamp of started meeting
	Timestamp *int64 `json:"timestamp,omitempty"`
}

func NewMeetingDTO(meetingID string, members []UserDTO, timestamp *int64) *MeetingDTO {
	return &MeetingDTO{
		MeetingID: meetingID,
		Members:   members,
		Timestamp: timestamp,
	}
}
