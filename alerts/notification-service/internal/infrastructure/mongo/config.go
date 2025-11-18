package mongo

import (
	"fmt"
	"net/url"
	"notification-serivce/internal/config"
	"os"
)

type MongoConfig struct {
	Host string
	Port string
	User string
	Pass string

	URI string

	DbName     string
	Collection string
}

func NewMongoConfig() *MongoConfig {
	uri := os.Getenv(config.MONGO_URI)
	host := os.Getenv(config.MONGO_HOST)
	port := os.Getenv(config.MONGO_PORT)
	user := os.Getenv(config.MONGO_USER)
	pass := os.Getenv(config.MONGO_PASS)

	dbname := os.Getenv(config.MONGO_DB_NAME)
	collection := os.Getenv(config.MONGO_COLLECTION)

	if _, err := url.Parse(uri); err != nil || uri == "" {

		if host == "" || port == "" || user == "" || pass == "" {
			return nil
		}

		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
	}

	return &MongoConfig{
		URI:        uri,
		Host:       host,
		Pass:       pass,
		User:       user,
		Port:       port,
		DbName:     dbname,
		Collection: collection,
	}
}
