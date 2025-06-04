package dto

type VoiceSessionDTO struct {
	ID        string   `json:"id"`
	Duration  int64    `json:"duration"`          // in seconds
	StartedAt string   `json:"startedAt"`         // ISO 8601 format
	EndedAt   *string  `json:"endedAt,omitempty"` // ISO 8601 format, optional
	MemberIDs []string `json:"memberIDs"`
	AudioURL  *string  `json:"audioURL,omitempty"`
}

func (v *VoiceSessionDTO) IsFinished() bool {
	return v.EndedAt != nil
}
