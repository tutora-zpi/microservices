package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoClient(rootCtx context.Context, mongoConfig MongoConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(rootCtx, mongoConfig.Timeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoConfig.URL()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	defer func() {
		if err != nil {
			if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
				log.Printf("Failed to disconnect MongoDB client after error: %v", disconnectErr)
			}
		}
	}()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("failed to ping MongoDB due to context error: %w", ctx.Err())
		}
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func Close(rootCtx context.Context, client *mongo.Client, timeout time.Duration) error {
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(rootCtx, timeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("failed to disconnect MongoDB due to context error: %w", ctx.Err())
		}
		return fmt.Errorf("failed to disconnect MongoDB: %w", err)
	}

	log.Println("Successfully disconnected from MongoDB")
	return nil
}
