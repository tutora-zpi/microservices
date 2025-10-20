package repository

import (
	"context"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/models"
	"meeting-scheduler-service/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type plannedMeetingsRepoImpl struct {
	db *gorm.DB
}

// CancelMeeting implements repository.PlannedMeetingsRepository.
func (p *plannedMeetingsRepoImpl) CancelMeeting(ctx context.Context, id int) error {
	tx := p.db.WithContext(ctx).
		Where("id = ? AND is_processed = false", id).
		Delete(&models.PlannedMeeting{})

	if tx.Error != nil {
		log.Printf("Error during cancellation: %v\n", tx.Error)
		return fmt.Errorf("failed to cancel meeting")
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("no unprocessed meeting found with id=%d", id)
	}

	return nil
}

// GetPlannedMeetings implements repository.PlannedMeetingsRepository.
func (p *plannedMeetingsRepoImpl) GetPlannedMeetings(
	ctx context.Context,
	req dto.FetchPlannedMeetingsDTO,
) ([]dto.PlannedMeetingDTO, error) {
	var meetings []models.PlannedMeeting

	db := p.db.WithContext(ctx).
		Order("start_date DESC").
		Limit(req.Limit)

	if req.LastPlannedDate != nil {
		db = db.Where("start_date < ?", req.LastPlannedDate)
	}

	db = db.Where("classId = ?", req.ClassID)

	if err := db.Find(&meetings).Error; err != nil {
		return nil, err
	}

	results := make([]dto.PlannedMeetingDTO, len(meetings))
	for i := len(meetings) - 1; i >= 0; i-- {
		log.Println(meetings[i].StartDate)
		dto, _ := meetings[i].ToPlannedMeetingDTO()
		if dto != nil {
			results[len(meetings)-i-1] = *dto
		}
	}

	return results, nil
}

// CanStartAnotherMeeting implements repository.PlannedMeetingsRepository.
func (p *plannedMeetingsRepoImpl) CanStartAnotherMeeting(ctx context.Context, meeting dto.PlanMeetingDTO) bool {
	var model models.PlannedMeeting

	startOfDay := meeting.StartDate.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := p.db.WithContext(ctx).
		Where("class_id = ? AND start_date >= ? AND start_date < ? AND finish_date > ?",
			meeting.ClassID, startOfDay, endOfDay, meeting.StartDate).
		First(&model).Error

	return err != nil
}

// CreatePlannedMeetings implements repository.PlannedMeetingsRepository.
func (p *plannedMeetingsRepoImpl) CreatePlannedMeetings(ctx context.Context, meeting dto.PlanMeetingDTO) (*dto.PlannedMeetingDTO, error) {
	toInsert, err := models.NewPlannedMeeting(meeting)
	if err != nil {
		return nil, fmt.Errorf("invalid format of user data")
	}
	if err := p.db.WithContext(ctx).Save(&toInsert).Error; err != nil {
		return nil, fmt.Errorf("failed to save planned meeting: %v", err)
	}

	return toInsert.ToPlannedMeetingDTO()
}

// ProcessPlannedMeetings implements repository.PlannedMeetingsRepository.
func (p *plannedMeetingsRepoImpl) ProcessPlannedMeetings(
	ctx context.Context,
	start, before time.Time,
) ([]dto.PlanMeetingDTO, error) {

	var meetings []models.PlannedMeeting

	tx := p.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.
		Where("is_processed = false AND start_date BETWEEN ? AND ?", start, before).
		Find(&meetings).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get planned meetings: %w", err)
	}

	if len(meetings) == 0 {
		tx.Rollback()
		return nil, nil
	}

	ids := make([]uint, len(meetings))
	for i, m := range meetings {
		ids[i] = m.ID
	}
	if err := tx.Model(&models.PlannedMeeting{}).
		Where("id IN ?", ids).
		Update("is_processed", true).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to mark as processed: %w", err)
	}

	results := make([]dto.PlanMeetingDTO, len(meetings))
	for i, m := range meetings {
		dto, err := m.ToPlanMeetingDTO()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to convert meeting: %w", err)
		}
		results[i] = *dto
	}

	return results, tx.Commit().Error
}

func NewPlannedMeetingsRepository(db *gorm.DB) repository.PlannedMeetingsRepository {
	return &plannedMeetingsRepoImpl{db: db}
}
