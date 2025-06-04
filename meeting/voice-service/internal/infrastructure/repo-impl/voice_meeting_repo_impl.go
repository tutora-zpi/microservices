package repoimpl

import (
	"time"
	"voice-service/internal/domain/dto"
	"voice-service/internal/domain/event"
	"voice-service/internal/domain/mapper"
	"voice-service/internal/domain/model"
	"voice-service/internal/domain/repository"
	"voice-service/internal/infrastructure/database"
)

type voiceMeetingRepositoryImpl struct {
	postgres database.Postgres
}

func NewVoiceMeetingRepository(db database.Postgres) repository.VoiceMeetingRepository {
	return &voiceMeetingRepositoryImpl{
		postgres: db,
	}
}

// FinishSession implements repository.VoiceMeetingRepository.
func (v *voiceMeetingRepositoryImpl) FinishSession(endedAt time.Time, id, audioURL string) (*dto.VoiceSessionDTO, error) {
	session := &model.VoiceSession{
		ID:       id,
		AudioURL: &audioURL,
		EndedAt:  &endedAt,
	}

	if err := v.postgres.Orm().Save(session).Error; err != nil {
		return nil, err
	}

	dto := mapper.NewVoiceSessionDTO(*session)

	return &dto, nil
}

func (v *voiceMeetingRepositoryImpl) CreateSession(event event.MeetingStartedEvent) (*dto.VoiceSessionDTO, error) {
	session := model.NewVoiceSession(event)

	if err := v.postgres.Orm().Create(session).Error; err != nil {
		return nil, err
	}

	dto := mapper.NewVoiceSessionDTO(*session)

	return &dto, nil
}
