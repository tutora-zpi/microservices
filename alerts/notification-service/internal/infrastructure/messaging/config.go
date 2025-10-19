package messaging

import (
	"fmt"
	"log"
	"notification-serivce/internal/config"
	"os"
)

type RabbitConfig struct {
	User string
	Pass string
	Port string
	Host string

	URL string

	// How many times to retry connection
	Retries int

	NotificationExchange string
}

func NewRabbitMQConfig() *RabbitConfig {
	user := os.Getenv(config.RABBITMQ_DEFAULT_USER)
	pass := os.Getenv(config.RABBITMQ_DEFAULT_PASS)
	host := os.Getenv(config.RABBITMQ_HOST)
	port := os.Getenv(config.RABBITMQ_PORT)
	url := os.Getenv(config.RABBITMQ_URL)
	exchange := os.Getenv(config.NOTIFICATION_EXCHANGE)

	if url == "" {
		if user == "" || pass == "" || host == "" || port == "" {
			return nil
		}
		url = fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
	}

	log.Println(url)

	return &RabbitConfig{
		User:                 user,
		Pass:                 pass,
		Host:                 host,
		Port:                 port,
		URL:                  url,
		Retries:              3,
		NotificationExchange: exchange,
	}
}
