package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"voice-service/internal/infrastructure/config"
	"voice-service/internal/infrastructure/security"
)

const (
	auth   string = "Authorization"
	bearer string = "Bearer"
	id     string = "userID"
)

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte(os.Getenv(config.JWT_SECRET))

		header := r.Header.Get(auth)
		tokenStr := strings.Split(header, bearer)[0]

		userID, err := security.DecodeJWT(tokenStr, secret)

		if err != nil {
			log.Printf("Unauthorized access: %s\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), id, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
