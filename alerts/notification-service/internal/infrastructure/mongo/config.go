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

	Uri string

	DbName     string
	Collection string
}

func (m *MongoConfig) URL() string {
	uri := os.Getenv(config.MONGO_URI)
	if _, err := url.Parse(uri); err != nil || uri == "" {
		host := os.Getenv(config.MONGO_HOST)
		port := os.Getenv(config.MONGO_PORT)
		user := os.Getenv(config.MONGO_USER)
		pass := os.Getenv(config.MONGO_PASS)

		if host == "" || port == "" || user == "" || pass == "" {
			return ""
		}

		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
	}

	return uri
}
