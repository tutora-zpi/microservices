package messaging

import (
	"log"
	"voice-service/internal/app/interfaces"
	"voice-service/internal/domain/event"
	"voice-service/internal/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func connect(connstr string, retries int) (*amqp.Connection, *amqp.Channel, error) {
	var err error
	var conn *amqp.Connection
	var ch *amqp.Channel

	for attempts := retries; attempts > 0; attempts-- {
		conn, err = amqp.Dial(connstr)

		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v. Retrying (%d/%d)...", err, retries-attempts+1, retries)
			continue
		}

		ch, err = conn.Channel()

		if err != nil {
			log.Printf("Failed to open a channel: %v. Retrying (%d/%d)...", err, retries-attempts+1, retries)
			conn.Close()
			continue
		}

		break
	}

	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil

}

// Close implements Broker.
func (r *Rabbit) Close() {
	if err := r.channel.Close(); err != nil {
		log.Println("Failed to close channel")
	}

	if err := r.connection.Close(); err != nil {
		log.Println("Failed to close connection")
	}
}

// ConsumeEvent implements Broker.
func (r *Rabbit) Consume(event.EventWrapper) error {
	panic("unimplemented")
}

// Publish implements Broker.
func (r *Rabbit) Publish(event event.EventWrapper) error {
	panic("unimplemented")
}

func NewRabbitBroker(cfg config.RabbitConfig) interfaces.Broker {

	conn, ch, err := connect(cfg.Connstr, cfg.Retries)

	if err != nil {
		panic("Failed to connect to RabbitMQ, check your config or RabbitMQ")
	}

	return &Rabbit{
		connection: conn,
		channel:    ch,
		queueName:  cfg.Queue,
	}
}
