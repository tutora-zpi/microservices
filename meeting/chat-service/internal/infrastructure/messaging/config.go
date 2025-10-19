package messaging

import (
	"chat-service/internal/config"
	"fmt"
	"net/url"
	"os"
	"time"
)

type RabbitConfig struct {
	User string
	Pass string
	Port string
	Host string
	URL  string

	ChatExchange    string
	MeetingExchange string

	Timeout      time.Duration
	ExchangeType string
}

func NewRabbitConfig() *RabbitConfig {
	return &RabbitConfig{
		Pass: os.Getenv(config.RABBITMQ_DEFAULT_PASS),
		User: os.Getenv(config.RABBITMQ_DEFAULT_USER),
		Host: os.Getenv(config.RABBITMQ_HOST),
		Port: os.Getenv(config.RABBITMQ_PORT),
		URL:  os.Getenv(config.RABBITMQ_URL),

		ChatExchange:    os.Getenv(config.CHAT_EXCHANGE),
		MeetingExchange: os.Getenv(config.MEETING_EXCHANGE),

		Timeout:      time.Second * 5,
		ExchangeType: "fanout",
	}
}

func (r *RabbitConfig) GetURL() string {
	_, err := url.Parse(r.URL)
	if err == nil && r.URL != "" {
		return r.URL
	}

	if r.User == "" || r.Pass == "" || r.Host == "" || r.Port == "" {
		return ""
	}

	return fmt.Sprintf("amqp://%s:%s@%s:%s", r.User, r.Pass, r.Host, r.Port)
}
