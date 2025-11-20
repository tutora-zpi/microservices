package request

import "strconv"

const DEFAULT_LIMIT = 10

type FetchSessions struct {
	MeetingID     string
	LastFetchedID *string
	Limit         int64
}

func NewFetchSessions(meetingID, lastFetchedID, limit string) *FetchSessions {
	var result FetchSessions

	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		limitNum = DEFAULT_LIMIT
	}

	result.Limit = int64(limitNum)

	if lastFetchedID != "" {
		result.LastFetchedID = &lastFetchedID
	}

	result.MeetingID = meetingID

	return &result
}
