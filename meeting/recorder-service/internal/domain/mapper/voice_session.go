package mapper

import (
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/model"
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
