package handlers

import (
	"fmt"
	"net/http"
	"time"
)

const (
	UserIDContextKey         = "userID"
	DefaultHeartbeatInterval = 15 * time.Second
)

func ExtractClientID(r *http.Request) (string, error) {
	val := r.Context().Value(UserIDContextKey)
	clientID, ok := val.(string)
	if !ok || clientID == "" {
		return "", fmt.Errorf("missing client ID in context")
	}
	return clientID, nil
}
