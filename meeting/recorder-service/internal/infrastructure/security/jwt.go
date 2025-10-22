package security

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func DecodeJWT(tokenStr string, secret []byte) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims or signature")
	}

	id, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid sub claim")
	}

	if _, err := uuid.Parse(id); err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	return id, nil
}
