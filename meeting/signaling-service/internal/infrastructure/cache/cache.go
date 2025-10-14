package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenService interface {
	DoesTokenExists(ctx context.Context, token string) bool
	SaveToken(ctx context.Context, token string, ttl time.Duration) error
}

type tokenServiceImpl struct {
	client   *redis.Client
	tokenKey func(suffix string) string
}

// SaveToken implements TokenService.
func (t *tokenServiceImpl) SaveToken(ctx context.Context, token string, ttl time.Duration) error {
	key := t.tokenKey(token)
	value := fmt.Sprint(ttl)

	if t.DoesTokenExists(ctx, token) {
		log.Println("Token had been saved before")
		return nil
	}

	_, err := t.client.SetEx(ctx, key, value, ttl).Result()
	if err != nil {
		log.Printf("Something went wrong: %v", err)
		return fmt.Errorf("failed to set %s under %s", value, key)
	}

	return nil
}

// DoesTokenExists implements TokenService.
func (t *tokenServiceImpl) DoesTokenExists(ctx context.Context, token string) bool {
	key := t.tokenKey(token)

	existsNumber, err := t.client.Exists(ctx, key).Result()

	if err != nil {
		log.Printf("Not found token with: %s", key)
		return false
	}

	return existsNumber == 1
}

func NewTokenService(client *redis.Client) TokenService {
	return &tokenServiceImpl{
		client: client,
		tokenKey: func(suffix string) string {
			return fmt.Sprintf("token:%s", suffix)
		},
	}
}
