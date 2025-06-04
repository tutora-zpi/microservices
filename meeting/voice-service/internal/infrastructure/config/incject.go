package config

import (
	"voice-service/internal/app/interfaces"
	"voice-service/internal/domain/repository"
	"voice-service/internal/infrastructure/ws"
)

type Incjectable struct {
	Broker   interfaces.Broker
	Recorder interfaces.Recorder
	Repo     repository.VoiceMeetingRepository
	Gateway  ws.Gateway
}
