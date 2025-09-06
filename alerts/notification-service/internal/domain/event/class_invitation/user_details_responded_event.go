package classinvitation

import (
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/models"
	"reflect"
)

type UserDetails struct {
	dto.UserDTO
	Role string `json:"role"`
}

type UserDetailsRespondedEvent struct {
	ID       string      `json:"notificationId"`
	Sender   UserDetails `json:"sender"`
	Receiver UserDetails `json:"receiver"`
}

func (u *UserDetailsRespondedEvent) Name() string {
	return reflect.TypeOf(u).Elem().Name()
}

func (u *UserDetailsRespondedEvent) FieldsToUpdate() map[string]any {
	return map[string]any{
		"receiver": *models.NewUser(u.Receiver.ID, u.Receiver.FirstName, u.Receiver.LastName, u.Receiver.Role),
		"sender":   *models.NewUser(u.Sender.ID, u.Sender.FirstName, u.Sender.LastName, u.Sender.Role),
		"status":   enums.CREATED,
	}
}
