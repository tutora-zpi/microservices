package factoryimpl

import (
	"recorder-service/internal/app/interfaces/factory"
	irecorder "recorder-service/internal/domain/recorder"
	"recorder-service/internal/infrastructure/webrtc/recorder"
)

type recorderFactoryImpl struct {
}

// CreateNewRecorder implements factory.RecorderFactory.
func (r *recorderFactoryImpl) CreateNewRecorder() irecorder.Recorder {
	return recorder.NewRecorderClient()
}

func NewRecorderFactory() factory.RecorderFactory {
	return &recorderFactoryImpl{}
}
