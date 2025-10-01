package messaging

import (
	"fmt"
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/event"
	"meeting-scheduler-service/internal/infrastructure/config"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const fanout string = "fanout"

const RETRIES int = 3

var rabbitmqURL string

type RabbitMQBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel

	publishingCfg amqp.Publishing
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

// Publish implements Broker.
func (r *RabbitMQBroker) Publish(eventToSend event.Event, exchangeChannels ...string) error {
	for _, channel := range exchangeChannels {
		if err := r.reconnect(channel); err != nil {
			log.Printf("Failed to reconnect to rabbitmq: %v\n", err)
			return err
		}

		wrapper := event.EventWrapper{
			Pattern: eventToSend.Name(),
			Data:    eventToSend,
		}

		body, err := wrapper.ToJson()
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		err = r.channel.Publish(
			channel,
			wrapper.Pattern,
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

		log.Printf("Published event to exchange %s: %s", channel, string(body))
	}
	return nil
}

func NewRabbitBroker() interfaces.Broker {

	rabbitmqURL = buildConnectionString()

	conn, err := connect(rabbitmqURL)
	if err != nil {
		log.Panicln("Failed to connect to RabbitMQ, check your config or RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to open channel: %v\n", err)
	}

	exchangeName := os.Getenv(config.EVENT_EXCHANGE_QUEUE_NAME)
	notificationExchangeChannel := os.Getenv(config.NOTIFICATION_EXCHANGE_QUEUE_NAME)

	if err := declareExchange(ch, exchangeName, notificationExchangeChannel); err != nil {
		log.Panicf("Failed to declare exchange: %v\n", err)
	}

	log.Println("Successfully connected to RabbitMQ")

	return &RabbitMQBroker{
		connection: conn,
		channel:    ch,
	}
}

func (r *RabbitMQBroker) reconnect(exchangeChannels ...string) error {
	for _, channel := range exchangeChannels {
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

			if err := declareExchange(r.channel, channel); err != nil {
				return err
			}
		}
	}
	return nil
}

func declareExchange(ch *amqp.Channel, exchangeNames ...string) error {
	for _, exchangeName := range exchangeNames {
		if err := ch.ExchangeDeclare(
			exchangeName,
			fanout,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to declare exchange after reconnect: %w", err)
		}

	}
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
