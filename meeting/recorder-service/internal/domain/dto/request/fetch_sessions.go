package request

import "strconv"

const DEFAULT_LIMIT = 10

type FetchSessions struct {
	ClassID              string
	LastFetchedMeetingID *string
	Limit                int64
}

func NewFetchSessions(classID, lastFetchedMeetingID, limit string) *FetchSessions {
	var result FetchSessions

	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		limitNum = DEFAULT_LIMIT
	}

	result.Limit = int64(limitNum)

	if lastFetchedMeetingID != "" {
		result.LastFetchedMeetingID = &lastFetchedMeetingID
	}

	result.ClassID = classID

	return &result
}
