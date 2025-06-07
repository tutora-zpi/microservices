package messaging

import (
	"fmt"
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/event"
	"meeting-scheduler-service/internal/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel

	publishingCfg amqp.Publishing
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
func (r *RabbitMQBroker) Close() {
	if err := r.channel.Close(); err != nil {
		log.Println("Failed to close channel")
	}

	if err := r.connection.Close(); err != nil {
		log.Println("Failed to close connection")
	}
}

// ConsumeEvent implements Broker.
func (r *RabbitMQBroker) Consume(event.EventWrapper) error {
	panic("unimpl")
}

// func preprocess(messages <-chan amqp.Delivery) {
// 	for d := range messages {
// 		log.Printf("Received a message: %s", d.Body)
// 		dotCount := bytes.Count(d.Body, []byte("."))
// 		t := time.Duration(dotCount)
// 		time.Sleep(t * time.Second)
// 		log.Printf("Done")
// 	}
// }

// Publish implements Broker.
func (r *RabbitMQBroker) Publish(event event.EventWrapper) error {
	ch, err := r.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.EVENT_EXCHANGE_QUEUE_NAME,
		"fanout", // type
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	body, err := event.ToJson()
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = ch.Publish(
		config.EVENT_EXCHANGE_QUEUE_NAME, // exchange
		"",                               // routing key
		false,                            // mandatory
		false,                            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func NewRabbitBroker(cfg config.RabbitConfig) interfaces.Broker {

	conn, ch, err := connect(cfg.Connstr, cfg.Retries)

	if err != nil {
		panic("Failed to connect to RabbitMQ, check your config or RabbitMQ")
	}

	return &RabbitMQBroker{
		connection: conn,
		channel:    ch,
	}
}
