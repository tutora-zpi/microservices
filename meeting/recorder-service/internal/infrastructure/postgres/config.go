package postgres

import (
	"fmt"
	"os"
	"recorder-service/internal/config"
	"strconv"
	"time"
)

type PostgresConfig struct {
	User     string
	Password string

	Host string
	Port string

	URL string

	DatabaseName string
	SSLMode      bool

	Timeout time.Duration
}

func NewPostgresConfig(timeout time.Duration) *PostgresConfig {
	sslMode, err := strconv.ParseBool(os.Getenv(config.POSTGRES_SSLMODE))
	if err != nil {
		sslMode = false
	}

	sslModeStr := "disable"
	if sslMode {
		sslModeStr = "enable"
	}

	user := os.Getenv(config.POSTGRES_USER)
	pass := os.Getenv(config.POSTGRES_PASS)
	host := os.Getenv(config.POSTGRES_HOST)
	port := os.Getenv(config.POSTGRES_PORT)
	databaseName := os.Getenv(config.POSTGRES_DBNAME)

	url := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, pass, databaseName, port, sslModeStr,
	)

	return &PostgresConfig{
		User:         user,
		Password:     pass,
		Host:         host,
		Port:         port,
		URL:          url,
		DatabaseName: databaseName,
		SSLMode:      sslMode,
		Timeout:      timeout,
	}
}
