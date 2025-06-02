package repoimpl

import "voice-service/internal/domain/repository"

type voiceMeetingRepositoryImpl struct{}

func NewVoiceMeetingRepository() repository.VoiceMeetingRepository {
	return &voiceMeetingRepositoryImpl{}
}
