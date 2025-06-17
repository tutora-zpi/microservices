package mapper

import (
	"voice-service/internal/domain/dto"
	"voice-service/internal/domain/model"
)

func NewVoiceSessionDTO(model model.VoiceSession) dto.VoiceSessionDTO {
	ended := model.EndedAt.String()

	return dto.VoiceSessionDTO{
		ID:        model.ID,
		Duration:  int64(model.EndedAt.Sub(model.StartedAt).Seconds()),
		StartedAt: model.StartedAt.String(),
		EndedAt:   &ended,
		MemberIDs: model.MemberIDs,
		AudioURL:  model.AudioURL,
	}
}
