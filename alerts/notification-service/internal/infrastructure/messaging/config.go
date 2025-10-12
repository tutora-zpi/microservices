package messaging

import (
	"fmt"
	"net/url"
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
