package model

import (
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/event"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VoiceSessionMetadata struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	MeetingID string        `bson:"meetingId"`
	ClassID   string        `bson:"classId"`

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

func (v *VoiceSessionMetadata) DTO() *dto.VoiceSessionMetadataDTO {
	return &dto.VoiceSessionMetadataDTO{
		ID:        v.ID.Hex(),
		ClassID:   v.ClassID,
		MeetingID: v.MeetingID,
		EndedAt:   v.EndedAt,
		StartedAt: &v.StartedAt,
		Duration:  int64(v.EndedAt.Sub(v.StartedAt)),
		MemberIDs: v.MemberIDs,
		AudioName: v.MergedAudioName,
	}
}
