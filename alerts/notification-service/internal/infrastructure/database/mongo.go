package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	db *mongo.Collection
}

const DATABASE_NAME string = "notification_db"
const COLLECTION string = "notifications"

func Connect() (*Database, error) {

	uri := getConnectionURL()
	if uri == "" {
		return nil, fmt.Errorf("MongoDB credentials are missing")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("Failed to ping mongo db")
	}

	log.Println("Database pinged successfully!")

	return chooseDatabase(client)
}

func (d *Database) Close() error {
	log.Println("Closing connections...")

	if err := d.db.Database().Client().Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("Failed to close connection: %s", err.Error())
	}

	log.Println("MongoDB connection closed successfully")
	return nil
}

func (d *Database) GetCollection() *mongo.Collection {
	return d.db
}

func chooseDatabase(client *mongo.Client) (*Database, error) {
	dbName := os.Getenv("MONGO_DB_NAME")

	if dbName == "" {
		dbName = DATABASE_NAME
	}

	collection := os.Getenv("MONGO_COLLECTION")

	if collection == "" {
		collection = COLLECTION
	}

	mongoDB := client.Database(dbName).Collection(collection)

	return &Database{
		db: mongoDB,
	}, nil
}

func getConnectionURL() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		host := os.Getenv("MONGO_HOST")
		port := os.Getenv("MONGO_PORT")
		user := os.Getenv("MONGO_USER")
		pass := os.Getenv("MONGO_PASS")

		if host == "" || port == "" || user == "" || pass == "" {
			return ""
		}

		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
	}

	return uri
}
