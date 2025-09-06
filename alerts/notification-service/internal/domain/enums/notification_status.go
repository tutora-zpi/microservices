package enums

type NotificationStatus string

const (
	DELIVERED NotificationStatus = "delivered"
	CREATED   NotificationStatus = "created"
	PENDING   NotificationStatus = "pending"
)
