package repoimpl

import (
	"fmt"
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

// AppendAudioURL implements repository.VoiceMeetingRepository.
func (v *voiceMeetingRepositoryImpl) AppendAudioURL(id string, audioURL string) error {
	if err := v.postgres.Orm().Where("id = ?", id).Update("audio_url = ?", audioURL).Error; err != nil {
		return fmt.Errorf("could not update audio in %s", id)
	}
	return nil
}

func NewVoiceMeetingRepository(db database.Postgres) repository.VoiceMeetingRepository {
	return &voiceMeetingRepositoryImpl{
		postgres: db,
	}
}

// FinishSession implements repository.VoiceMeetingRepository.
func (v *voiceMeetingRepositoryImpl) FinishSession(id string) (*dto.VoiceSessionDTO, error) {
	t := time.Now()

	session := &model.VoiceSession{
		ID:      id,
		EndedAt: &t,
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
		return nil, fmt.Errorf("could not create session %s", err.Error())
	}

	dto := mapper.NewVoiceSessionDTO(*session)

	return &dto, nil
}
