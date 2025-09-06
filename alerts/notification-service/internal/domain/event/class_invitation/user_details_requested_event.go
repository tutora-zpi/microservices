package classinvitation

import (
	"notification-serivce/internal/domain/dto"
	"reflect"
)

type UserDetailsRequestedEvent struct {
	ID         string `json:"notificationId"`
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
}

func (u *UserDetailsRequestedEvent) Name() string {
	return reflect.TypeOf(u).Elem().Name()
}

func NewUserDetailsRequestedEvent(dto *dto.NotificationDTO) *UserDetailsRequestedEvent {
	return &UserDetailsRequestedEvent{
		ID:         dto.ID,
		SenderID:   dto.Sender.ID,
		ReceiverID: dto.Receiver.ID,
	}
}
