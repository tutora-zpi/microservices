package dto

// MeetingDTO represents meeting details returned in responses.
// swagger:model MeetingDTO
type MeetingDTO struct {
	// Meeting unique identifier
	MeetingID string `json:"meetingID"`
	// Members who participated in the meeting
	Members []UserDTO `json:"members"`
}
