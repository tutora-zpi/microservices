package messaging

import (
	"chat-service/internal/config"
	"fmt"
	"os"
	"time"
)

type RabbitConfig struct {
	User string
	Pass string
	Port string
	Host string

	URL string

	Exchanges    []string
	ExchangeType string
	ChatQueue    string
	FileQueue    string

	PoolSize int
	Timeout  time.Duration
}

func NewRabbitMQConfig(timeout time.Duration, poolSize int) *RabbitConfig {
	user := os.Getenv(config.RABBITMQ_DEFAULT_USER)
	pass := os.Getenv(config.RABBITMQ_DEFAULT_PASS)
	host := os.Getenv(config.RABBITMQ_HOST)
	port := os.Getenv(config.RABBITMQ_PORT)
	url := os.Getenv(config.RABBITMQ_URL)

	if url == "" {
		if user == "" || pass == "" || host == "" || port == "" {
			return nil
		}
		url = fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
	}

	return &RabbitConfig{
		User:         user,
		Pass:         pass,
		Host:         host,
		Port:         port,
		URL:          url,
		PoolSize:     poolSize,
		Timeout:      timeout,
		ExchangeType: "fanout",
		ChatQueue:    os.Getenv(config.CHAT_QUEUE),
		Exchanges: []string{
			os.Getenv(config.MEETING_EXCHANGE),
			os.Getenv(config.CHAT_EXCHANGE),
		},
		FileQueue: os.Getenv(config.FILE_QUEUE),
	}
}
