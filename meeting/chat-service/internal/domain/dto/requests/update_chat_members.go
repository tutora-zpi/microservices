package requests

import "github.com/go-playground/validator/v10"

type UpdateChatMembers struct {
	ChatID     string   `json:"chatId" validate:"required"`
	MembersIDs []string `json:"membersIds" validate:"required,min=1,dive,required"`
}

func (u *UpdateChatMembers) IsValid() error {
	v := validator.New()

	return v.Struct(u)
}
