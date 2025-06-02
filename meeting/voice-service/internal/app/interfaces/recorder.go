package interfaces

type Recorder interface {
	StartRecording(fileName string)
	StopRecording()
}
