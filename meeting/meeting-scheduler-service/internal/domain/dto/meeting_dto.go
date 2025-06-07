package dto

type MeetingDTO struct {
	MeetingID string    `json:"meetingID"`
	Members   []UserDTO `json:"members"`
}
