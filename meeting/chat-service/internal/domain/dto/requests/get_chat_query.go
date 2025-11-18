package requests

type GetChat struct {
	ID    string `json:"id" validate:"reqiured,uuid4"`
	Limit int    `json:"limit"`
}
