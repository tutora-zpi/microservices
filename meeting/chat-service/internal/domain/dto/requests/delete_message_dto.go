package requests

type DeleteMessage struct {
	ChatID    string `json:"chatId" validate:"required,uuid4"`
	MessageID string `json:"messageId" validate:"required,uuid4"`
}
