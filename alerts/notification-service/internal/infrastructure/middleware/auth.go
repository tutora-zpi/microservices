package middleware

import (
	"context"
	"log"
	"net/http"
	"notification-serivce/internal/infrastructure/security"
	"notification-serivce/internal/infrastructure/server"
	"notification-serivce/pkg"
	"strings"
)

const (
	auth   string = "Authorization"
	bearer string = "Bearer "
	id     string = "userID"
	token  string = "token"
)

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(auth)
		var tokenStr string = ""

		if header == "" {
			tokenStr = r.URL.Query().Get(token)
		} else {
			tokenStr = strings.TrimPrefix(header, bearer)
		}

		// if tokenStr == "" {
		// 	log.Println("Missing token")
		// 	server.NewResponse(w, pkg.Ptr("Missing JWT token"), http.StatusUnauthorized, nil)
		// 	return
		// }

		userID, err := security.DecodeJWT(tokenStr)

		if err != nil {
			log.Printf("Unauthorized access: %s\n", err)
			server.NewResponse(w, pkg.Ptr("Token expired or invalid"), http.StatusUnauthorized, nil)
			return
		}

		ctx := context.WithValue(r.Context(), id, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
