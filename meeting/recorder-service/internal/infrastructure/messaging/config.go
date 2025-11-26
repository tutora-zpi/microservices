package messaging

import (
	"fmt"
	"os"
	"recorder-service/internal/config"
	"time"
)

type RabbitConfig struct {
	User string
	Pass string
	Port string
	Host string

	URL string

	MeetingExchange string
	ExchangeType    string
	RecorderQueue   string

	PoolSize int
	Timeout  time.Duration

	PrefetchCount int
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
		User:            user,
		Pass:            pass,
		Host:            host,
		Port:            port,
		URL:             url,
		PoolSize:        poolSize,
		Timeout:         timeout,
		ExchangeType:    "fanout",
		MeetingExchange: os.Getenv(config.MEETING_EXCHANGE),
		RecorderQueue:   os.Getenv(config.RECORDER_QUEUE),
		PrefetchCount:   10,
	}
}
