package dto

import (
	"reflect"
	"strconv"
)

const (
	DEFAULT_LIMIT int = 10
)

type FetchNotificationsDTO struct {
	ReceiverID         string
	Limit              int
	LastNotificationID *string
}

func (f *FetchNotificationsDTO) Name() string {
	return reflect.TypeOf(f).Elem().Name()
}

func NewFetchNotificationsDTO(limit, lastNotificationID, receiverID string) *FetchNotificationsDTO {
	q := &FetchNotificationsDTO{}

	q.ReceiverID = receiverID

	limitNum, err := strconv.Atoi(limit)

	if err != nil || limitNum > 20 || limitNum < 1 {
		q.Limit = DEFAULT_LIMIT
	} else {
		q.Limit = limitNum
	}

	if lastNotificationID != "" {
		q.LastNotificationID = &lastNotificationID
	}

	return q
}
