package interfaces

type EventHandler interface {
	Handle(body []byte) error
}
