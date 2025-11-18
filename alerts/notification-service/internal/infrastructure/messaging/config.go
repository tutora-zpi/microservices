package messaging

import (
	"fmt"
	"notification-serivce/internal/config"
	"os"
	"time"
)

type RabbitConfig struct {
	User string
	Pass string
	Port string
	Host string

	URL string

	NotificationQueue string

	PoolSize int
	Timeout  time.Duration

	ExchangeType string

	Exchanges []string
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
		User:              user,
		Pass:              pass,
		Host:              host,
		Port:              port,
		URL:               url,
		NotificationQueue: os.Getenv(config.NOTIFICATION_QUEUE),
		PoolSize:          poolSize,
		ExchangeType:      "fanout",
		Timeout:           timeout,
		Exchanges: []string{
			os.Getenv(config.MEETING_EXCHANGE),
			os.Getenv(config.CLASS_EXCHANGE),
		},
	}
}
