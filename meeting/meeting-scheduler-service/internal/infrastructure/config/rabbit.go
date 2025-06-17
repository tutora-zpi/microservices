package config

import (
	"log"
	"strings"
)

type RabbitConfig struct {
	// Connection string to RabbitMQ server
	Connstr string

	// How many times to retry connection
	Retries int
}

func NewRabbitConfig(connstr string, retries int) RabbitConfig {
	if connstr == "" {
		panic("RabbitMQ connection string cannot be empty")
	}

	log.Printf("Connecting to broker on: %s", strings.Split(connstr, "@")[1])

	return RabbitConfig{
		Connstr: connstr,
		Retries: retries,
	}
}
