package messaging

import (
	"fmt"
	"log"
	"notification-serivce/internal/config"
	"notification-serivce/internal/domain/event"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

const RETRIES int = 3

var rabbitmqURL string

func init() {
	rabbitmqURL = buildConnectionString()
}

func connect(connStr string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for attempts := range RETRIES {
		conn, err = amqp.Dial(connStr)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v. Retrying (%d/%d)...", err, attempts+1, RETRIES)
			time.Sleep(3 * time.Second)
			continue
		}
		return conn, nil
	}

	return nil, fmt.Errorf("could not connect to RabbitMQ: %w", err)
}

func NewRabbitBroker() *RabbitMQBroker {

	conn, err := connect(rabbitmqURL)
	if err != nil {
		log.Panicln("Failed to connect to RabbitMQ, check your config or RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to open channel: %v\n", err)
	}

	exchangeName := os.Getenv(config.EVENT_EXCHANGE_QUEUE_NAME)
	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		log.Panicf("Failed to declare exchange: %v\n", err)
	}

	return &RabbitMQBroker{
		connection: conn,
		channel:    ch,
		exchange:   exchangeName,
	}
}

func (r *RabbitMQBroker) Close() {
	if err := r.channel.Close(); err != nil {
		log.Println("Failed to close channel:", err)
	}
	if err := r.connection.Close(); err != nil {
		log.Println("Failed to close connection:", err)
	}
}

func (r *RabbitMQBroker) Publish(event event.EventWrapper) error {
	if r.channel == nil || r.connection.IsClosed() {
		connStr := buildConnectionString()

		conn, err := connect(connStr)
		if err != nil {
			return fmt.Errorf("reconnect failed: %w", err)
		}
		ch, err := conn.Channel()
		if err != nil {
			return fmt.Errorf("failed to open channel after reconnect: %w", err)
		}
		r.connection = conn
		r.channel = ch

		if err := ch.ExchangeDeclare(
			r.exchange,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to declare exchange after reconnect: %w", err)
		}
	}

	body, err := event.ToJson()
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = r.channel.Publish(
		r.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published event to exchange %s: %s", r.exchange, string(body))
	return nil
}

func buildConnectionString() string {
	url := os.Getenv(config.RABBITMQ_URL)

	if url == "" {
		pass := os.Getenv(config.RABBITMQ_DEFAULT_PASS)
		user := os.Getenv(config.RABBITMQ_DEFAULT_USER)
		port := os.Getenv(config.RABBITMQ_PORT)
		host := os.Getenv(config.RABBITMQ_HOST)

		if pass == "" || user == "" || host == "" || port == "" {
			return ""
		}

		url = fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
	}

	return url
}
