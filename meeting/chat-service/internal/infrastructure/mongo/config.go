package mongo

import (
	"chat-service/internal/config"
	"fmt"
	"net/url"
	"os"
)

type MongoConfig struct {
	Host string
	Port string
	Pass string
	User string
	Uri  string

	Database string

	ReactionCollection string
	MessagesCollection string
	ChatCollection     string
}

func NewMongoConfig() *MongoConfig {
	return &MongoConfig{
		Host: os.Getenv(config.MONGO_HOST),
		Port: os.Getenv(config.MONGO_PORT),
		User: os.Getenv(config.MONGO_USER),
		Pass: os.Getenv(config.MONGO_PASS),

		Uri: os.Getenv(config.MONGO_URI),

		Database:           os.Getenv(config.MONGO_DATABASE),
		ReactionCollection: os.Getenv(config.MONGO_REACTION_COLLECTION),
		ChatCollection:     os.Getenv(config.CHAT_EXCHANGE),
		MessagesCollection: os.Getenv(config.MONGO_MESSAGES_COLLECTION),
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
