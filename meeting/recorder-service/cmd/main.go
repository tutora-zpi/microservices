package main

import (
	"os"
	"recorder-service/internal/config"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		_ = godotenv.Load(".env.local")
	}

	// security.FetchSignKey()
}

func main() {

}
