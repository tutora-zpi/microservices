package model

import (
	"path"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/event"
	"time"
)

type VoiceSessionMetadata struct {
	MeetingID string `bson:"meetingId"`
	ClassID   string `bson:"classId"`

	StartedAt time.Time  `bson:"startedAt"`
	EndedAt   *time.Time `bson:"endedAt,omitempty"`
	MemberIDs []string   `bson:"memberIds"`

	MergedAudioName *string `bson:"audioName,omitempty"`
}

func NewVoiceSession(event event.MeetingStartedEvent) *VoiceSessionMetadata {

	ids := make([]string, len(event.Members))

	for i, member := range event.Members {
		ids[i] = member.ID
	}

	return &VoiceSessionMetadata{
		MeetingID: event.MeetingID,
		ClassID:   event.ClassID,
		StartedAt: event.StartedTime,
		EndedAt:   &event.FinishTime,
		MemberIDs: ids,
	}
}

func (v *VoiceSessionMetadata) GetAudioURL() string {
	path := path.Join("api", "v1", "recordings", v.MeetingID)
	return path
}

func (v *VoiceSessionMetadata) DTO() *dto.VoiceSessionMetadataDTO {
	audioURL := v.GetAudioURL()

	return &dto.VoiceSessionMetadataDTO{
		ClassID:   v.ClassID,
		MeetingID: v.MeetingID,
		EndedAt:   v.EndedAt,
		StartedAt: &v.StartedAt,
		Duration:  int64(v.EndedAt.Sub(v.StartedAt)),
		MemberIDs: v.MemberIDs,
		AudioURL:  &audioURL,
	}
}
