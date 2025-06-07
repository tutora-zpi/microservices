package middleware

import (
	"context"
	"log"
	"meeting-scheduler-service/internal/infrastructure/config"
	"meeting-scheduler-service/internal/infrastructure/security"
	"net/http"
	"os"
	"strings"
)

const (
	auth   string = "Authorization"
	bearer string = "Bearer "
	id     string = "userID"
)

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte(os.Getenv(config.JWT_SECRET))

		header := r.Header.Get(auth)
		tokenStr := strings.TrimPrefix(header, bearer)

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
