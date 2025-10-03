package middleware

import (
	"context"
	"encoding/json"
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/infrastructure/server"
	"net/http"
)

type key string

const dtoKey key = "dto"

var dtoRegistry = map[string]func() any{
	"/api/v1/meeting/start": func() any { return &dto.StartMeetingDTO{} },
	"/api/v1/meeting/end":   func() any { return &dto.EndMeetingDTO{} },
}

func Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		constructor, ok := dtoRegistry[r.URL.Path]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		dtoInstance := constructor()
		if err := json.NewDecoder(r.Body).Decode(dtoInstance); err != nil {
			server.NewResponse(w, "Failed to decode body", http.StatusBadRequest, nil)
			return
		}

		switch v := dtoInstance.(type) {
		case *dto.StartMeetingDTO:
			if err := v.IsValid(); err != nil {
				server.NewResponse(w, "Validation error: "+err.Error(), http.StatusBadRequest, nil)
				return
			}
		case *dto.EndMeetingDTO:
			if err := v.IsValid(); err != nil {
				server.NewResponse(w, "Validation error: "+err.Error(), http.StatusBadRequest, nil)
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
