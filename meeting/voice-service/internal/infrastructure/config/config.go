package config

import (
	"gorm.io/gorm"
)

type PostgresConfig struct {
	// Connection string to PostgreSQL server
	Connstr string

	// How many times to retry connection
	Retries int

	// Models to migrate
	ModelsToMigrate []any

	// Optional GORM configuration
	Options *gorm.Config
}

func NewPostgresConfig(connstr string, retries int, options *gorm.Config, models ...any) PostgresConfig {
	if connstr == "" {
		panic("Postgres connection string cannot be empty")
	}

	return PostgresConfig{
		Connstr:         connstr,
		Retries:         retries,
		Options:         options,
		ModelsToMigrate: models,
	}
}

type RabbitConfig struct {
	// Connection string to RabbitMQ server
	Connstr string

	// Name of queue to consume messages from or publish
	Queue string

	// How many times to retry connection
	Retries int
}

func NewRabbitConfig(connstr string, queue string, retries int) RabbitConfig {
	if connstr == "" {
		panic("RabbitMQ connection string cannot be empty")
	}
	if queue == "" {
		panic("RabbitMQ queue name cannot be empty")
	}

	return RabbitConfig{
		Connstr: connstr,
		Queue:   queue,
		Retries: retries,
	}
}
