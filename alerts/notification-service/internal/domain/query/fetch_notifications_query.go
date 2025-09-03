package query

import "reflect"

const (
	DEFAULT_LIMIT int = 10
)

type FetchNotificationsQuery struct {
	// person who requests it
	ReceiverID         string
	Limit              int
	LastNotificationID *string
}

func (f *FetchNotificationsQuery) Name() string {
	return reflect.TypeOf(f).Elem().Name()
}

func NewFetchNotificationsQuery(limit, lastNotificationID, receiverID string) Query {
	q := &FetchNotificationsQuery{}

	q.ReceiverID = receiverID

	if limit == "" || limit > "20" {
		q.Limit = DEFAULT_LIMIT
	}

	if lastNotificationID != "" {
		q.LastNotificationID = &lastNotificationID
	}

	return q
}
