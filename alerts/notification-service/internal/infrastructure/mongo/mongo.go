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

func NewMongoClient(ctx context.Context, config MongoConfig, timeout time.Duration) (*mongo.Client, error) {
	if config.DbName == "" || config.Collection == "" {
		return nil, fmt.Errorf("set database name, collection name")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(config.URI))
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

func Close(ctx context.Context, client *mongo.Client, timeout time.Duration) {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		if ctx.Err() != nil {
			log.Printf("failed to disconnect MongoDB due to context error: %v", ctx.Err())
			return
		}
		log.Printf("failed to disconnect MongoDB: %v", err)
	}

	log.Println("Successfully disconnected from MongoDB")

}
