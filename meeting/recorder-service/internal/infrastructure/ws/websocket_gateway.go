package ws

import (
	"net/http"
)

type Gateway interface {
	Handle(w http.ResponseWriter, r *http.Request)
	Broadcast(to string, message []byte) error
}
