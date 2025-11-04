package mongo

import (
	"fmt"
	"net/url"
	"os"
	"recorder-service/internal/config"
	"time"
)

type MongoConfig struct {
	Host string
	Port string
	Pass string
	User string
	Uri  string

	Database string

	SessionMetadataCollection string

	Timeout time.Duration
}

func NewMongoConfig(timeout time.Duration) *MongoConfig {
	return &MongoConfig{
		Host: os.Getenv(config.MONGO_HOST),
		Port: os.Getenv(config.MONGO_PORT),
		User: os.Getenv(config.MONGO_USER),
		Pass: os.Getenv(config.MONGO_PASS),

		Uri: os.Getenv(config.MONGO_URI),

		Database:                  os.Getenv(config.MONGO_DATABASE),
		SessionMetadataCollection: os.Getenv(config.MONGO_METADATA_COLLECTION),
		Timeout:                   timeout,
	}
}

func (m *MongoConfig) URL() string {
	uri := m.Uri
	if _, err := url.Parse(uri); err != nil || uri == "" {
		if m.Host == "" || m.Port == "" || m.User == "" || m.Pass == "" {
			return ""
		}

		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", m.User, m.Pass, m.Host, m.Port)
	}

	return uri
}
