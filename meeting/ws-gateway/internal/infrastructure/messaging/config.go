package messaging

import (
	"fmt"
	"net/url"
	"os"
	"ws-gateway/internal/config"
)

type RabbitConfig struct {
	User string
	Pass string
	Port string
	Host string

	URL string

	ExchangeNames []string
}

func NewRabbitMQConfig() *RabbitConfig {
	return &RabbitConfig{
		User: os.Getenv(config.RABBITMQ_DEFAULT_USER),
		Pass: os.Getenv(config.RABBITMQ_DEFAULT_PASS),
		URL:  os.Getenv(config.RABBITMQ_URL),
		Host: os.Getenv(config.RABBITMQ_HOST),
		Port: os.Getenv(config.RABBITMQ_PORT),

		ExchangeNames: []string{
			os.Getenv(config.CHAT_EXCHANGE),
			os.Getenv(config.BOARD_EXCHANGE),
		},
	}
}

func (r *RabbitConfig) RabbitMQURL() string {
	_, err := url.Parse(r.URL)
	if err == nil && r.URL != "" {
		return r.URL
	}

	if r.User == "" || r.Pass == "" || r.Host == "" || r.Port == "" {
		return ""
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%s", r.User, r.Pass, r.Host, r.Port)

}
