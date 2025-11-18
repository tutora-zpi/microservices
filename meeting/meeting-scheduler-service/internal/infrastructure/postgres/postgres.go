package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres(ctx context.Context, postgresConfig PostgresConfig) (*gorm.DB, error) {
	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	ctx, cancel := context.WithTimeout(ctx, postgresConfig.Timeout)
	defer cancel()

	db, err := gorm.Open(postgres.Open(postgresConfig.URL), &gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database")
	}

	sql, err := db.DB()
	if err != nil {
		Close(ctx, db, postgresConfig)
		return nil, fmt.Errorf("no database")
	}

	if err := sql.PingContext(ctx); err != nil {
		Close(ctx, db, postgresConfig)
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("ping timeout exceeded (%s): %w", postgresConfig.Timeout, err)
		}
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	if err := db.AutoMigrate(&models.PlannedMeeting{}); err != nil {
		Close(ctx, db, postgresConfig)
		return nil, fmt.Errorf("failed to migrate models: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	return db, nil
}

func Close(ctx context.Context, client *gorm.DB, postgresConfig PostgresConfig) {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, postgresConfig.Timeout)
	defer cancel()

	var errCh chan error = make(chan error, 1)

	go func() {
		defer close(errCh)

		sql, err := client.DB()
		if err != nil {
			errCh <- fmt.Errorf("failed to get DB connection: %w", err)
			return
		}

		if err := sql.Close(); err != nil {
			errCh <- fmt.Errorf("failed to close db: %w", err)
			return
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("Closing PostreSQL took too much time")
	case err := <-errCh:
		if err != nil {
			log.Printf("Failed to close PostgreSQL connection: %v", err)
		} else {
			log.Println("Successfully closed PostgreSQL connection")
		}
	}
}
