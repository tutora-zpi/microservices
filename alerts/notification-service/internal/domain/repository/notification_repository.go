package repository

type NotificationRepository interface {
	Save()
	MarkAsRead()
}
