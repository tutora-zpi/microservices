package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/domain/broker"
	"signaling-service/internal/domain/event"
	"signaling-service/internal/infrastructure/bus"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	fanout = "fanout"
)

type RabbitMQBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	dispatcher bus.Dispachable

	config RabbitConfig
}

func NewRabbitBroker(rabbitMQConfig RabbitConfig, dispatcher bus.Dispachable) (interfaces.Broker, error) {
	conn, err := connect(rabbitMQConfig.Retries, rabbitMQConfig.RabbitMQURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	if err := declareExchanges(ch, rabbitMQConfig.ExchangeNames...); err != nil {
		return nil, fmt.Errorf("failed to declare exchanges: %v", err)
	}

	log.Println("Successfully connected to RabbitMQ")

	return &RabbitMQBroker{
		connection: conn,
		channel:    ch,
		dispatcher: dispatcher,
		config:     rabbitMQConfig,
	}, nil
}

func (r *RabbitMQBroker) Close() {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.connection != nil {
		_ = r.connection.Close()
	}

	log.Println("RabbitMQ connection successfully closed")
}

func (r *RabbitMQBroker) Publish(ctx context.Context, ev event.Event, dest broker.Destination) error {
	if err := r.ensureConnected(dest.Exchange); err != nil {
		return fmt.Errorf("failed to ensure connection: %w", err)
	}

	wrapper := event.EventWrapper{
		Pattern: ev.Name(),
		Data:    ev,
	}

	body, err := json.Marshal(wrapper)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	pub := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Timestamp:   time.Now(),
	}

	var exchange string
	var routingKey string

	switch {
	case dest.Exchange != "":
		exchange = dest.Exchange
		routingKey = dest.RoutingKey
	case dest.Queue != "":
		exchange = ""
		routingKey = dest.Queue
	default:
		return fmt.Errorf("no destination specified")
	}

	if err := r.channel.PublishWithContext(ctx, exchange, routingKey, false, false, pub); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published event [%s] to %s:%s", wrapper.Pattern, exchange, routingKey)
	return nil
}

func (r *RabbitMQBroker) PublishMultiple(ctx context.Context, ev event.Event, destinations ...broker.Destination) error {
	for _, dest := range destinations {
		if err := r.Publish(ctx, ev, dest); err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitMQBroker) ensureConnected(exchangeNames ...string) error {
	if r.connection == nil || r.connection.IsClosed() {
		conn, err := connect(r.config.Retries, r.config.RabbitMQURL())
		if err != nil {
			return err
		}

		ch, err := conn.Channel()
		if err != nil {
			return err
		}

		r.connection = conn
		r.channel = ch

		if len(exchangeNames) > 0 {
			if err := declareExchanges(r.channel, exchangeNames...); err != nil {
				return err
			}
		}
	}
	return nil
}

func declareExchanges(ch *amqp.Channel, exchangeNames ...string) error {
	for _, name := range exchangeNames {
		if name == "" {
			continue
		}
		if err := ch.ExchangeDeclare(name, fanout, true, false, false, false, nil); err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", name, err)
		}
	}
	return nil
}

func connect(retries int, url string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for attempt := 1; attempt <= retries; attempt++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			return conn, nil
		}
		log.Printf("Failed to connect to RabbitMQ: %v (retry %d/%d)", err, attempt, retries)
		time.Sleep(3 * time.Second)
	}
	return nil, fmt.Errorf("could not connect to RabbitMQ after retries: %w", err)
}
