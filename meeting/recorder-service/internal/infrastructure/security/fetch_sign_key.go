package security

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"recorder-service/internal/config"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"token_type"`
	ExipresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	BotID       string
}

func FetchSignKey() {
	var err error
	jwksURL := os.Getenv(config.JWKS_URL)
	if jwksURL == "" {
		log.Fatalln("jwks url is empty")
	}

	JWKS, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			log.Printf("JWKS refresh error: %v", err)
		},
		RefreshUnknownKID: true,
	})

	if err != nil {
		log.Fatalf("Failed to get JWKS from %s: %v", jwksURL, err)
	}
	log.Println("Successfully initialized JWKS")
}

func FetchToken(ctx context.Context) (*TokenResponse, error) {
	urlPath := os.Getenv(config.TOKEN_URL)
	if urlPath == "" {
		return nil, fmt.Errorf("TOKEN_URL is empty")
	}

	clientSecret := os.Getenv(config.OAUTH_BOT_CLIENT_SECRET)
	clientID := os.Getenv(config.OAUTH_BOT_CLIENT_ID)

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(respBody, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	id, err := DecodeJWT(tokenResponse.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve id from access token: %w", err)
	}

	tokenResponse.BotID = id

	return &tokenResponse, nil
}
