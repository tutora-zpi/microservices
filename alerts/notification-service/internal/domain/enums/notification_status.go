package enums

type NotificationStatus string

const (
	Delivered NotificationStatus = "delivered"
	Created   NotificationStatus = "created"
	Pending   NotificationStatus = "pending"
)
