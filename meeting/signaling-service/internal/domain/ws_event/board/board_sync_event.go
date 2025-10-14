package board

import "github.com/go-playground/validator/v10"

type BoardSyncEvent struct {
	Data map[string]any `json:"data"`
}

func (b *BoardSyncEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(b)
}

func (b *BoardSyncEvent) Type() string {
	return "board:sync"
}

func (u *BoardSyncEvent) Name() string {
	return u.Type()
}
