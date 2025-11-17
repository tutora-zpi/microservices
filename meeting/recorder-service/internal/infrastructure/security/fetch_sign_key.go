package security

import (
	"context"
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

	clientSecret := os.Getenv(config.CLIENT_SECRET)
	clientID := os.Getenv(config.CLIENT_ID)

	clientCredentials := url.Values{}
	clientCredentials.Set("grant_type", "client_credentials")
	clientCredentials.Set("client_id", clientID)
	clientCredentials.Set("client_secret", clientSecret)

	body := strings.NewReader(clientCredentials.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, body)
	log.Printf("Making REQUEST: %v", *req)

	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil || len(respBody) < 1 {
		return nil, fmt.Errorf("failed to read response body")
	}

	var tokenResponse TokenResponse

	if err := json.Unmarshal(respBody, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response")
	}
	id, err := DecodeJWT(tokenResponse.AccessToken)

	if err != nil {
		return nil, fmt.Errorf("failed to decode access token")
	}

	tokenResponse.BotID = id

	return &tokenResponse, nil
}
