package dto

type PlannedMeetingDTO struct {
	// Identifier of planned meeting (UUIDv4)
	ID int `json:"id"`

	PlanMeetingDTO
}
