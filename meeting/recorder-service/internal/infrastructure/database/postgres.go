package database

import (
	"fmt"
	"log"
	"recorder-service/internal/infrastructure/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres interface {
	Connect(config config.PostgresConfig) error

	Migrate(models []any) error

	Close() error

	Orm() *gorm.DB
}

type postgresDatabaseImpl struct {
	gorm *gorm.DB
}

func (p *postgresDatabaseImpl) Migrate(models []any) error {
	if len(models) == 0 {
		log.Println("No models provided for migration")
		return nil
	}

	if err := p.gorm.AutoMigrate(models...); err != nil {
		log.Printf("Failed to migrate models: %v", err)
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	log.Println("Models migrated successfully")
	return nil
}

func (p *postgresDatabaseImpl) Orm() *gorm.DB {
	if p.gorm == nil {
		log.Println("GORM DB is not initialized")
		return nil
	}

	return p.gorm
}

func NewPostgres(cfg config.PostgresConfig) Postgres {
	db := &postgresDatabaseImpl{}
	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil
	}
	log.Println("PostgreSQL connection established")

	if err := db.Migrate(cfg.ModelsToMigrate); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
		return nil
	}
	log.Println("PostgreSQL models migrated successfully")

	return db
}

func (p *postgresDatabaseImpl) Close() error {
	sqlDB, err := p.gorm.DB()
	if err != nil {
		log.Printf("Failed to get sql.DB: %v", err)
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Failed to close DB: %v", err)
		return fmt.Errorf("failed to close DB: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}

func (p *postgresDatabaseImpl) Connect(config config.PostgresConfig) error {
	var err error
	var db *gorm.DB

	if config.Options == nil {
		config.Options = &gorm.Config{}
	}

	for attempts := config.Retries; attempts > 0; attempts-- {
		db, err = gorm.Open(postgres.Open(config.Connstr), config.Options)
		if err == nil {
			p.gorm = db
			log.Println("Connected to PostgreSQL successfully")
			return nil
		}
		log.Printf("Failed to connect to PostgreSQL, retries left: %d, error: %v\n", attempts-1, err)
	}

	return fmt.Errorf("failed to connect to PostgreSQL after %d retries: %w", config.Retries, err)
}
