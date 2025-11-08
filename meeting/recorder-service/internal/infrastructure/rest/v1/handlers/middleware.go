package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"recorder-service/internal/infrastructure/security"
	"recorder-service/internal/infrastructure/server"
	"recorder-service/pkg"
	"strings"
)

const (
	Auth         string = "Authorization"
	BearerPrefix string = "Bearer "
	ID           string = "userID"
	Token        string = "token"
)

func findToken(r *http.Request) (string, error) {
	var token string

	cookie, err := r.Cookie(Token)
	if err == nil {
		token = cookie.Value
		return token, nil
	}
	log.Println("Not found token in cookie going to find in query")

	token = r.URL.Query().Get(Token)
	if token != "" {
		return token, nil
	}

	log.Println("Not found token in query going to find in header")

	auth := r.Header.Get(Auth)
	if !strings.HasPrefix(auth, BearerPrefix) {
		return "", fmt.Errorf("no bearer prefix in header")
	}

	token = strings.TrimPrefix(auth, BearerPrefix)
	token = strings.TrimSpace(token)

	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

func (h *handlerImpl) IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := findToken(r)

		if err != nil {
			log.Printf("Not found token: %s\n", err)
			server.NewResponse(w, pkg.Ptr("Unauthorized access"), http.StatusUnauthorized, nil)
			return
		}

		userID, err := security.DecodeJWT(tokenStr)

		if err != nil {
			log.Printf("Unauthorized access: %s\n", err)
			server.NewResponse(w, pkg.Ptr("Unauthorized access"), http.StatusUnauthorized, nil)
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, ID, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
