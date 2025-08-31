package security

import (
	"log"
	"notification-serivce/internal/config"
	"os"
	"time"

	"github.com/MicahParks/keyfunc"
)

func FetchSignKey() {
	var err error
	jwksURL := os.Getenv(config.JWKS_URL)
	if jwksURL == "" {
		log.Panicln("jwks url is empty")
	}

	JWKS, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			log.Printf("JWKS refresh error: %v", err)
		},
		RefreshUnknownKID: true,
	})

	if err != nil {
		log.Panicf("Failed to get JWKS from %s: %v\n", jwksURL, err)
	}
	log.Println("Successfully initialized JWKS")
}
