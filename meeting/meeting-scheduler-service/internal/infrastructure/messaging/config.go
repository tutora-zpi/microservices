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

	ExchangeNames []string
}

func (this *RabbitConfig) RabbitMQURL() string {
	_, err := url.Parse(this.URL)
	if err == nil {
		return this.URL
	}

	if this.User == "" || this.Pass == "" || this.Host == "" || this.Port == "" {
		return ""
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%s", this.User, this.Pass, this.Host, this.Port)

}
