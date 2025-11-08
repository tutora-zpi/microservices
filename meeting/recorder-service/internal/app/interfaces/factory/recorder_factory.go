package factory

import "recorder-service/internal/domain/recorder"

type RecorderFactory interface {
	CreateNewRecorder() recorder.Recorder
}
