package interfaces

type Recorder interface {
	StartRecording(meetingID string) error
	StopRecording(meetingID string) (string, error)
}
