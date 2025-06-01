package database

import (
	"fmt"
	"log"
	"voice-service/internal/domain/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	connstr string
	retries int
	options *gorm.Config
}

type Postgres interface {
	Connect(config PostgresConfig) error
	Migrate([]model.Model) error
	Close() error
}

type PostgresDb struct {
	gorm *gorm.DB
}

// Close implements Postgres.
func (p *PostgresDb) Close() error {
	sqlDB, err := p.gorm.DB()
	if err != nil {
		log.Fatalln(err)

		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	defer sqlDB.Close()

	return nil
}

// Connect implements Postgres.
func (p *PostgresDb) Connect(config PostgresConfig) error {
	var err error
	var db *gorm.DB

	for err != nil && config.retries > 0 {
		log.Printf("Failed to connect to PostgreSQL, retries left: %d\n", config.retries)
		config.retries--
		db, err = gorm.Open(postgres.Open(config.connstr), config.options)
	}

	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	p.gorm = db
	log.Println("Connected to PostgreSQL successfully")
	return nil
}
