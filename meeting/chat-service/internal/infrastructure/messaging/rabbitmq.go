package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/broker"
	"chat-service/internal/domain/event"
	"chat-service/internal/infrastructure/bus"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	dispatcher bus.Dispachable
	config     RabbitConfig
	mu         sync.Mutex
}

func NewRabbitBroker(rabbitMQConfig RabbitConfig, dispatcher bus.Dispachable) (interfaces.Broker, error) {
	if rabbitMQConfig.Timeout == 0 {
		rabbitMQConfig.Timeout = 5 * time.Second
	}
	if rabbitMQConfig.ExchangeType == "" {
		rabbitMQConfig.ExchangeType = "fanout"
	}

	conn, err := connect(context.Background(), rabbitMQConfig.GetURL(), rabbitMQConfig.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if err := declareExchanges(ch, rabbitMQConfig.ExchangeType, rabbitMQConfig.ChatExchange, rabbitMQConfig.MeetingExchange); err != nil {
		conn.Close()
		ch.Close()
		return nil, fmt.Errorf("failed to declare exchanges: %w", err)
	}

	log.Printf("Successfully connected to RabbitMQ")

	return &RabbitMQBroker{
		conn:       conn,
		ch:         ch,
		dispatcher: dispatcher,
		config:     rabbitMQConfig,
	}, nil
}

func (r *RabbitMQBroker) Publish(ctx context.Context, ev event.Event, dest broker.Destination) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.ensureConnected(ctx, dest.Exchange); err != nil {
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
		Timestamp:   time.Now().UTC(),
	}

	var exchange, routingKey string
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

	if err := r.ch.PublishWithContext(ctx, exchange, routingKey, false, false, pub); err != nil {
		return fmt.Errorf("failed to publish message to %s:%s: %w", exchange, routingKey, err)
	}

	log.Printf("Published event [%s] to %s:%s", wrapper.Pattern, exchange, routingKey)
	return nil
}

func (r *RabbitMQBroker) PublishMultiple(ctx context.Context, ev event.Event, destinations ...broker.Destination) error {
	var firstErr error
	for _, dest := range destinations {
		if err := r.Publish(ctx, ev, dest); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			log.Printf("Failed to publish to %v: %v", dest, err)
		}
	}
	return firstErr
}

func (r *RabbitMQBroker) Consume(ctx context.Context, exchange string) error {
	if err := r.ensureConnected(ctx, exchange); err != nil {
		return fmt.Errorf("failed to ensure connection: %w", err)
	}

	q, err := r.ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	for _, p := range r.dispatcher.AvailablePatterns() {
		if err := r.ch.QueueBind(q.Name, p, exchange, false, nil); err != nil {
			return fmt.Errorf("failed to bind queue to %s with pattern %s: %w", exchange, p, err)
		}
	}

	msgs, err := r.ch.ConsumeWithContext(ctx, q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	log.Printf("Started consumer on queue %s for exchange %s", q.Name, exchange)

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer stopped due to context cancellation")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				log.Println("Consumer channel closed")
				return fmt.Errorf("consumer channel closed")
			}
			if len(msg.Body) == 0 {
				if err := msg.Ack(false); err != nil {
					log.Printf("Failed to ack empty message: %v", err)
				}
				continue
			}

			var wrapper event.EventWrapper
			pattern, data, err := wrapper.DecodedEventWrapper(msg.Body)
			if err != nil {
				if err := msg.Nack(false, false); err != nil {
					log.Printf("Failed to nack message: %v", err)
				}
				continue
			}

			if err := r.dispatcher.HandleEvent(ctx, pattern, data); err != nil {
				log.Printf("Failed to handle event %s: %v", pattern, err)
				if err := msg.Nack(false, true); err != nil {
					log.Printf("Failed to nack message: %v", err)
				}
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}
	}

}

func (r *RabbitMQBroker) ensureConnected(ctx context.Context, exchangeNames ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn == nil || r.conn.IsClosed() {
		if r.ch != nil {
			if err := r.ch.Close(); err != nil {
				log.Printf("Failed to close existing channel: %v", err)
			}
			r.ch = nil
		}
		if r.conn != nil {
			if err := r.conn.Close(); err != nil {
				log.Printf("Failed to close existing connection: %v", err)
			}
			r.conn = nil
		}

		conn, err := connect(ctx, r.config.GetURL(), r.config.Timeout)
		if err != nil {
			return fmt.Errorf("failed to reconnect to RabbitMQ: %w", err)
		}

		ch, err := conn.Channel()
		if err != nil {
			conn.Close()
			return fmt.Errorf("failed to open channel: %w", err)
		}

		r.conn = conn
		r.ch = ch

		if len(exchangeNames) > 0 {
			if err := declareExchanges(r.ch, r.config.ExchangeType, exchangeNames...); err != nil {
				return fmt.Errorf("failed to declare exchanges: %w", err)
			}
		}
	}
	return nil
}

func declareExchanges(ch *amqp.Channel, exchangeType string, exchangeNames ...string) error {
	for _, name := range exchangeNames {
		if name == "" {
			continue
		}
		if err := ch.ExchangeDeclare(
			name,
			exchangeType,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", name, err)
		}

		log.Printf("%s exchange declared successfully", name)
	}
	return nil
}

func connect(ctx context.Context, url string, timeout time.Duration) (*amqp.Connection, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to connect, timeout reached")
	default:
		if url == "" {
			return nil, fmt.Errorf("rabbitmq connection string is empty")
		}
		conn, err := amqp.Dial(url)
		if err != nil {
			return nil, fmt.Errorf("could not connect to RabbitMQ: %w", err)
		}
		ch, err := conn.Channel()
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to open channel to verify connection: %w", err)
		}
		ch.Close()

		return conn, nil
	}
}
