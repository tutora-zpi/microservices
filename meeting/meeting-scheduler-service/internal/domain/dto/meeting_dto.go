package dto

// MeetingDTO represents meeting details returned in responses.
// swagger:model MeetingDTO
type MeetingDTO struct {
	// Meeting unique identifier
	MeetingID string `json:"meetingId"`
	// Members who participated in the meeting
	Members []UserDTO `json:"members,omitempty"`
	// Timestamp of started meeting
	Timestamp *int64 `json:"timestamp,omitempty"`
	// Meetings title
	Title string `json:"title"`
}

func NewMeetingDTO(meetingID string, members []UserDTO, timestamp *int64, title string) *MeetingDTO {
	return &MeetingDTO{
		MeetingID: meetingID,
		Members:   members,
		Timestamp: timestamp,
		Title:     title,
	}
}
