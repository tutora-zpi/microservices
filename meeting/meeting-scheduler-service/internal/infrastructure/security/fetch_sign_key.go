package security

import (
	"log"
	"meeting-scheduler-service/internal/config"
	"os"
	"time"

	"github.com/MicahParks/keyfunc"
)

func FetchSignKey() {
	var err error
	jwksURL := os.Getenv(config.JWKS_URL)
	if jwksURL == "" {
		log.Fatalln("jwks url is empty")
		return
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
