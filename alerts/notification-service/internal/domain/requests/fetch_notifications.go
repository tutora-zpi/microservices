package requests

import (
	"reflect"
	"strconv"
)

const (
	DEFAULT_LIMIT int = 10
)

type FetchNotificationsRequest struct {
	// person who requests it
	ReceiverID         string
	Limit              int
	LastNotificationID *string
}

func (f *FetchNotificationsRequest) Name() string {
	return reflect.TypeOf(f).Elem().Name()
}

func NewFetchNotificationsRequest(limit, lastNotificationID, receiverID string) *FetchNotificationsRequest {
	q := &FetchNotificationsRequest{}

	q.ReceiverID = receiverID

	limitNum, err := strconv.Atoi(limit)

	if err != nil || limitNum > 20 {
		q.Limit = DEFAULT_LIMIT
	} else {
		q.Limit = limitNum
	}

	if lastNotificationID != "" {
		q.LastNotificationID = &lastNotificationID
	}

	return q
}
