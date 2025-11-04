package handlers

import (
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/validator"
	"chat-service/internal/infrastructure/server"
	"chat-service/pkg"
	"context"
	"encoding/json"
	"net/http"
)

type key string

const dtoKey key = "dto"

var dtoRegistry = map[string]func() any{
	"/api/v1/chats/general":        func() any { return &requests.CreateGeneralChat{} },
	"/api/v1/chats/update-members": func() any { return &requests.UpdateChatMembers{} },
}

func ValidateJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		constructor, ok := dtoRegistry[r.URL.Path]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		dtoInstance := constructor()
		if err := json.NewDecoder(r.Body).Decode(dtoInstance); err != nil {
			server.NewResponse(w, pkg.Ptr("Failed to decode body"), http.StatusBadRequest, nil)
			return
		}

		if validatable, ok := dtoInstance.(validator.Validable); ok {
			if err := validatable.IsValid(); err != nil {
				server.NewResponse(w, pkg.Ptr("Invalid body error: "+err.Error()), http.StatusBadRequest, nil)
				return
			}
		}

		ctx := context.WithValue(r.Context(), dtoKey, dtoInstance)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetDTO(r *http.Request) any {
	return r.Context().Value(dtoKey)
}
