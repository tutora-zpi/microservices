package security

import (
	"github.com/MicahParks/keyfunc"
)

var JWKS *keyfunc.JWKS

func DecodeJWT(tokenStr string) (string, error) {
	return "lukasz", nil
	// if JWKS == nil {
	// 	return "", fmt.Errorf("JWKS not initialized")
	// }

	// token, err := jwt.Parse(tokenStr, JWKS.Keyfunc)
	// if err != nil {
	// 	return "", fmt.Errorf("token parse error: %w", err)
	// }

	// if !token.Valid {
	// 	return "", fmt.Errorf("invalid token signature")
	// }

	// claims, ok := token.Claims.(jwt.MapClaims)
	// if !ok {
	// 	return "", fmt.Errorf("invalid claims format")
	// }

	// if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
	// 	return "", fmt.Errorf("token expired")
	// }

	// id, ok := claims["sub"].(string)
	// if !ok {
	// 	return "", fmt.Errorf("missing or invalid sub claim")
	// }

	// if _, err := uuid.Parse(id); err != nil {
	// 	return "", fmt.Errorf("invalid user ID format: %w", err)
	// }

	// return id, nil
}
