package dto

import (
	"fmt"
	"strconv"
	"time"
)

const (
	DEFAULT_LIMIT = 10
)

type FetchPlannedMeetingsDTO struct {
	ClassID         string     `json:"classId" validate:"required,uuid4"`
	Limit           int        `json:"limit" validate:"required"`
	LastPlannedDate *time.Time `json:"lastPlannedDate,omitempty"`
}

func NewFetchPlannedMeetingsDTO(classID, limit, lastStartTimestamp string) (*FetchPlannedMeetingsDTO, error) {
	var err error
	var date time.Time
	result := FetchPlannedMeetingsDTO{
		LastPlannedDate: nil,
		ClassID:         classID,
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		limitNum = DEFAULT_LIMIT
	}

	result.Limit = limitNum

	if len(lastStartTimestamp) > 0 {
		timestamp, err := strconv.Atoi(lastStartTimestamp)
		if err != nil {
			return nil, fmt.Errorf("date must be int")
		}

		date = time.Unix(int64(timestamp), 0).UTC()

		result.LastPlannedDate = &date
	}

	date = date.UTC().Truncate(time.Minute)

	return &result, nil
}
