package dto

import (
	"time"
)

type PlanMeetingDTO struct {
	StartMeetingDTO

	// Start Date - date and time when meeting starts, use toISOString to cast your date
	// required: true
	StartDate time.Time `json:"startDate" validate:"required" example:"2025-10-10T12:36:05+02:00"`
}

func (dto *PlanMeetingDTO) IsValid() error {
	if err := dto.StartMeetingDTO.IsValid(); err != nil {
		return err
	}

	// v := validator.New()
	// if err := v.Struct(dto); err != nil {
	// 	return err
	// }

	// dto.ConvertTimeToUTC()

	// if time.Until(dto.StartDate) < time.Minute*5 {
	// 	return fmt.Errorf("difference between now and start dates must be at least 5 minutes")
	// }

	// sub := dto.FinishDate.Sub(dto.StartDate)
	// if sub < 0 {
	// 	return fmt.Errorf("finish date is before start date")
	// }

	// if sub > time.Hour {
	// 	return fmt.Errorf("maximum length of meeting is 1 hour")
	// }

	return nil
}

func (dto *PlanMeetingDTO) ConvertTimeToUTC() {
	dto.StartDate = dto.StartDate.UTC().Truncate(time.Minute)
	dto.FinishDate = dto.FinishDate.UTC().Truncate(time.Minute)
}

func (dto *PlanMeetingDTO) GetDate() string {
	date := dto.StartDate.Format(time.DateOnly)
	return date
}
