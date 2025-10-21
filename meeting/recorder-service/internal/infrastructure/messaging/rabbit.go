package messaging

import (
	"context"
	"fmt"
	"log"
	"sync"

	"recorder-service/internal/app/interfaces"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connect(connstr string, retries int) (*amqp.Connection, *amqp.Channel, error) {
	var err error
	var conn *amqp.Connection
	var ch *amqp.Channel

	for attempts := retries; attempts > 0; attempts-- {
		conn, err = amqp.Dial(connstr)
		if err != nil {
			log.Printf("Failed to connect: %v. Retrying (%d/%d)...", err, retries-attempts+1, retries)
			continue
		}

		ch, err = conn.Channel()
		if err != nil {
			log.Printf("Failed to open channel: %v. Retrying (%d/%d)...", err, retries-attempts+1, retries)
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

type RabbitMQBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	mu         sync.Mutex
}

func (r *RabbitMQBroker) Close() {
	if err := r.channel.Close(); err != nil {
		log.Println("Failed to close channel:", err)
	}
	if err := r.connection.Close(); err != nil {
		log.Println("Failed to close connection:", err)
	}
}

func (r *RabbitMQBroker) Publish(event event.EventWrapper, opts config.PublishOptions) error {
	r.mu.Lock()
	ch, err := r.connection.Channel()
	r.mu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		opts.ExchangeName,
		opts.ExchangeType,
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
		opts.ExchangeName,
		opts.RoutingKey,
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

	return nil
}

func (r *RabbitMQBroker) Consume(ctx context.Context, options config.ConsumeOptions, dispacher interfaces.Dispatcher) error {
	r.mu.Lock()
	ch, err := r.connection.Channel()
	r.mu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	if options.ExchangeName != "" {
		err = ch.ExchangeDeclare(
			options.ExchangeName,
			options.ExchangeType,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare exchange: %w", err)
		}
	}

	q, err := ch.QueueDeclare(
		options.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if options.ExchangeName != "" {
		err = ch.QueueBind(
			q.Name,
			options.RoutingKey,
			options.ExchangeName,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go preprocess(ctx, msgs, dispacher.Dispatch)

	return nil
}

func preprocess(ctx context.Context, msgs <-chan amqp.Delivery, executor func(string, []byte) error) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer shutting down via context")
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}

			pattern := msg.RoutingKey

			if err := executor(pattern, msg.Body); err != nil {
				log.Printf("Handler error for pattern '%s': %v", pattern, err)
			}
		}
	}
}

func NewRabbitBroker(cfg config.RabbitConfig) interfaces.Broker {
	conn, ch, err := connect(cfg.Connstr, cfg.Retries)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	return &RabbitMQBroker{
		connection: conn,
		channel:    ch,
	}
}
