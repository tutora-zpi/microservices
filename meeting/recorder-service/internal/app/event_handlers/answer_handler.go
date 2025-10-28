package eventhandlers

import (
	"context"
	"recorder-service/internal/app/interfaces/handler"
)

type answerHandler struct {
}

// Handle implements interfaces.EventHandler.
func (a *answerHandler) Handle(ctx context.Context, body []byte) error {
	panic("unimplemented")
}

func NewAnswerHandler() handler.EventHandler {
	return &answerHandler{}
}
