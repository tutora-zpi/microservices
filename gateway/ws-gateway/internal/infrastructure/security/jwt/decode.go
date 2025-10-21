package security

import (
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var JWKS *keyfunc.JWKS

const (
	exp string = "exp"
	sub string = "sub"
)

func DecodeJWT(tokenStr string) (id string, ttl time.Duration, err error) {
	if JWKS == nil {
		return "", 0, fmt.Errorf("JWKS not initialized")
	}

	token, err := jwt.Parse(tokenStr, JWKS.Keyfunc)
	if err != nil {
		return "", 0, fmt.Errorf("token parse error: %w", err)
	}

	if !token.Valid {
		return "", 0, fmt.Errorf("invalid token signature")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, fmt.Errorf("invalid claims format")
	}

	exp, ok := claims[exp].(float64)

	if ok && int64(exp) < time.Now().UTC().Unix() {
		return "", 0, fmt.Errorf("token expired")
	}

	id, ok = claims[sub].(string)
	if !ok {
		return "", 0, fmt.Errorf("missing or invalid sub claim")
	}

	if _, err := uuid.Parse(id); err != nil {
		return "", 0, fmt.Errorf("invalid user ID format: %w", err)
	}

	ttl = time.Until(time.Unix(int64(exp), 0).UTC())

	return id, ttl, nil
}
