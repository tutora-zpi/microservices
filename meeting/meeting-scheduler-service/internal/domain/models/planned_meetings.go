package models

import (
	"encoding/json"
	"meeting-scheduler-service/internal/domain/dto"

	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PlannedMeeting struct {
	gorm.Model

	StartDate  time.Time `gorm:"index:idx_meeting_processed_start,priority:2;not null"`
	FinishDate time.Time `gorm:"not null"`

	Title       string         `gorm:"not null"`
	ClassID     string         `gorm:"not null"`
	Members     datatypes.JSON `gorm:"type:jsonb"`
	IsProcessed bool           `gorm:"index:idx_meeting_processed_start,priority:1;default:false"`
}

func (p *PlannedMeeting) ToPlannedMeetingDTO() (*dto.PlannedMeetingDTO, error) {
	planMeeting, err := p.ToPlanMeetingDTO()

	if err != nil {
		return nil, err
	}

	return &dto.PlannedMeetingDTO{
		ID:             int(p.ID),
		PlanMeetingDTO: *planMeeting,
	}, nil
}

func (p *PlannedMeeting) ToPlanMeetingDTO() (*dto.PlanMeetingDTO, error) {
	var members []dto.UserDTO
	if err := json.Unmarshal(p.Members, &members); err != nil {
		return nil, err
	}

	return &dto.PlanMeetingDTO{
		StartMeetingDTO: dto.StartMeetingDTO{
			ClassID:    p.ClassID,
			Title:      p.Title,
			Members:    members,
			FinishDate: p.FinishDate,
		},
		StartDate: p.StartDate,
	}, nil
}

func NewPlannedMeeting(dto dto.PlanMeetingDTO) (*PlannedMeeting, error) {
	membersJSON, err := json.Marshal(dto.Members)
	if err != nil {
		return nil, err
	}

	return &PlannedMeeting{
		StartDate:  dto.StartDate,
		FinishDate: dto.FinishDate,
		ClassID:    dto.ClassID,
		Title:      dto.Title,
		Members:    datatypes.JSON(membersJSON),
	}, nil
}
