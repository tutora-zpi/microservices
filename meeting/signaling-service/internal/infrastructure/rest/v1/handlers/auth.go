package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/infrastructure/security"
	"strings"
)

const (
	auth   string = "Authorization"
	bearer string = "Bearer "
	id     string = "userID"
	token  string = "token"
)

func (h *handlers) IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(auth)
		tokenStr := strings.TrimPrefix(header, bearer)

		userID, ttl, err := security.DecodeJWT(tokenStr)

		if err != nil {
			log.Printf("Unauthorized access: %s\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		if err := h.tokenService.SaveToken(ctx, tokenStr, ttl); err != nil {
			log.Printf("Failed to save token: %v", err)
		}

		ctx = context.WithValue(ctx, id, userID)
		ctx = context.WithValue(ctx, token, tokenStr)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *handlers) WithAuth(
	handler func(ctx context.Context, eventType string, msg []byte, client interfaces.Client) error,
) func(ctx context.Context, eventType string, msg []byte, client interfaces.Client) error {

	return func(ctx context.Context, eventType string, msg []byte, client interfaces.Client) error {
		tokenValue := ctx.Value(token)
		token, ok := tokenValue.(string)
		if !ok || !h.tokenService.DoesTokenExists(ctx, token) {
			client.GetConnection().Close()
			h.hub.RemoveGlobalMember(client)
			return fmt.Errorf("unauthorized")
		}

		return handler(ctx, eventType, msg, client)
	}
}
