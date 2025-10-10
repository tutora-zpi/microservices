package dto

type FetchPlannedMeetings struct {
	ClassID              string `json:"classId" validate:"required,uuid4"`
	Limit                int    `json:"limit" validate:"required"`
	LastPlannedMeetingID string `json:"lastPlannedMeetingId"`
}
