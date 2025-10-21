package usecase

import (
	"encoding/json"
	"log"
	"recorder-service/internal/app/interfaces"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/repository"
)

type RecordingMeeting struct {
	recorder interfaces.Recorder
	repo     repository.VoiceMeetingRepository
}

// Exec implements interfaces.UseCaseHandler.
func (r *RecordingMeeting) Exec(body []byte) error {
	var dest dto.RecordMessageDTO
	var path string
	err := json.Unmarshal(body, &dest)
	if err != nil {
		return err
	}

	switch dest.Command {
	case dto.Start:
		err = r.recorder.StartRecording(dest.MeetingID)
		if err != nil {
			log.Printf("Error during starting recording %s\n", err.Error())
		}
	case dto.Stop:
		path, err = r.recorder.StopRecording(dest.MeetingID)
		if err != nil {
			log.Printf("Error during stopping recording %s\n", err.Error())
		}

		err = r.repo.AppendAudioURL(dest.MeetingID, path) //. meeting always will be finished in future
	default:
		// possible early return from unmarshal
		log.Printf("Unknown %s command\n", dest.Command)
	}

	// generate notification about status of recording, extend record message dto in future

	return nil
}

func NewRecordingMeeting(recorder interfaces.Recorder, repo repository.VoiceMeetingRepository) interfaces.UseCaseHandler {
	return &RecordingMeeting{recorder: recorder, repo: repo}
}
