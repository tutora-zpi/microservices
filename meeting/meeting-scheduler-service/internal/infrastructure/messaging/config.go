package messaging

import (
	"fmt"
	"meeting-scheduler-service/internal/config"
	"os"
	"time"
)

type RabbitConfig struct {
	User string
	Pass string
	Port string
	Host string

	URL string

	NotificationExchange string
	MeetingExchange      string

	Timeout time.Duration
}

func NewRabbitMQConfig(timout time.Duration) *RabbitConfig {
	user := os.Getenv(config.RABBITMQ_DEFAULT_USER)
	pass := os.Getenv(config.RABBITMQ_DEFAULT_PASS)
	host := os.Getenv(config.RABBITMQ_HOST)
	port := os.Getenv(config.RABBITMQ_PORT)
	url := os.Getenv(config.RABBITMQ_URL)
	notificationExchange := os.Getenv(config.NOTIFICATION_EXCHANGE)
	meetingExchange := os.Getenv(config.MEETING_EXCHANGE)

	if url == "" {
		if user == "" || pass == "" || host == "" || port == "" {
			return nil
		}
		url = fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
	}

	return &RabbitConfig{
		User:                 user,
		Pass:                 pass,
		Host:                 host,
		Port:                 port,
		URL:                  url,
		NotificationExchange: notificationExchange,
		MeetingExchange:      meetingExchange,
		Timeout:              timout,
	}
}
