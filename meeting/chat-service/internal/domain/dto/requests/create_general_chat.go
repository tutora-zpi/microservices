package requests

import "github.com/go-playground/validator/v10"

type CreateGeneralChat struct {
	ClassID   string   `json:"classId" validate:"required,uuid4"`
	MemberIDs []string `json:"memberIds" validate:"required,dive"`
}

func (c *CreateGeneralChat) IsValid() error {
	v := validator.New()

	return v.Struct(c)
}
