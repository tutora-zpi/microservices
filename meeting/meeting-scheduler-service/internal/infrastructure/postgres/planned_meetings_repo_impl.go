package postgres

import (
	"context"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/models"
	"meeting-scheduler-service/internal/domain/repository"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type plannedMeetingsRepoImpl struct {
	db *gorm.DB
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

// Close implements repository.PlannedMeetingsRepository.
func (p *plannedMeetingsRepoImpl) Close() {
	sqlDB, _ := p.db.DB()
	sqlDB.Close()
}

// CreatePlannedMeetings implements repository.PlannedMeetingsRepository.
func (p *plannedMeetingsRepoImpl) CreatePlannedMeetings(ctx context.Context, meeting dto.PlanMeetingDTO) (*dto.PlanMeetingDTO, error) {
	toInsert, err := models.NewPlannedMeeting(meeting)
	if err != nil {
		return nil, fmt.Errorf("invalid format of user data")
	}
	if err := p.db.WithContext(ctx).Save(&toInsert).Error; err != nil {
		return nil, fmt.Errorf("failed to save planned meeting: %v", err)
	}

	return &meeting, nil
}

// ProcessPlannedMeetings implements repository.PlannedMeetingsRepository.
func (p *plannedMeetingsRepoImpl) ProcessPlannedMeetings(
	ctx context.Context,
	start time.Time,
	before time.Time,
) ([]dto.PlanMeetingDTO, error) {
	var meetings []models.PlannedMeeting
	var results []dto.PlanMeetingDTO

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("is_processed = false AND start_date >= ? AND start_date <= ?", start, before).
			Find(&meetings).Error; err != nil {
			return fmt.Errorf("failed to get planned meetings: %w", err)
		}

		results = make([]dto.PlanMeetingDTO, len(meetings))

		ids := make([]uint, len(meetings))
		for i, m := range meetings {
			m.IsProcessed = true

			dto, err := m.DTO()
			if err != nil {
				return fmt.Errorf("failed to convert meeting to DTO: %w", err)
			}
			results[i] = *dto
			ids[i] = m.ID
		}

		if len(ids) > 0 {
			if err := tx.Model(&models.PlannedMeeting{}).
				Where("id IN ?", ids).
				Update("is_processed", true).Error; err != nil {
				return fmt.Errorf("failed to mark meetings as processed: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func NewMeetingsRepository(postgresConfig PostgresConfig) (repository.PlannedMeetingsRepository, error) {
	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(postgres.Open(postgresConfig.ConnectionString()), &gormConfig)
	if err != nil {
		log.Printf("Failed to connect with PostgreSQL: %v", err)
		return nil, err
	}

	db.AutoMigrate(models.PlannedMeeting{})

	return &plannedMeetingsRepoImpl{db: db}, nil
}
