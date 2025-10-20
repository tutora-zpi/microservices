package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/broker"
	"notification-serivce/internal/domain/event"
	"notification-serivce/internal/infrastructure/bus"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	conn       *amqp.Connection
	chPool     chan *amqp.Channel
	dispatcher bus.Dispachable
	config     RabbitConfig
	connMu     sync.Mutex
}

func (r *RabbitMQBroker) Close() {
	r.connMu.Lock()
	defer r.connMu.Unlock()

	close(r.chPool)
	for ch := range r.chPool {
		if ch != nil {
			if err := ch.Close(); err != nil {
				log.Printf("Failed to close channel: %v", err)
			}
		}
	}

	if r.conn != nil && !r.conn.IsClosed() {
		if err := r.conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		} else {
			log.Println("RabbitMQ connection closed successfully")
		}
	}

	r.conn = nil
}

func NewRabbitBroker(rabbitMQConfig RabbitConfig, dispatcher bus.Dispachable) (interfaces.Broker, error) {
	if rabbitMQConfig.Timeout == 0 {
		rabbitMQConfig.Timeout = 5 * time.Second
	}

	conn, err := connect(context.Background(), rabbitMQConfig.URL, rabbitMQConfig.Timeout)
	if err != nil {
		return nil, err
	}

	broker := &RabbitMQBroker{
		conn:       conn,
		dispatcher: dispatcher,
		config:     rabbitMQConfig,
		chPool:     make(chan *amqp.Channel, rabbitMQConfig.PoolSize),
	}

	for i := 0; i < rabbitMQConfig.PoolSize; i++ {
		ch, err := conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to open channel: %w", err)
		}
		broker.chPool <- ch
	}

	firstCh := <-broker.chPool
	if err := declareExchanges(firstCh, rabbitMQConfig.ExchangeType, rabbitMQConfig.NotificationExchange); err != nil {
		return nil, err
	}
	broker.chPool <- firstCh

	log.Println("Successfully connected to RabittMQ")

	return broker, nil
}

func (r *RabbitMQBroker) Publish(ctx context.Context, ev event.Event, dest broker.Destination) error {
	ch := <-r.chPool
	defer func() { r.chPool <- ch }()

	r.connMu.Lock()
	if r.conn == nil || r.conn.IsClosed() {
		r.connMu.Unlock()
		return fmt.Errorf("connection is closed")
	}
	r.connMu.Unlock()

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
	if dest.Exchange != "" {
		exchange = dest.Exchange
		routingKey = dest.RoutingKey
	} else if dest.Queue != "" {
		exchange = ""
		routingKey = dest.Queue
	} else {
		return fmt.Errorf("no destination specified")
	}

	if err := ch.PublishWithContext(ctx, exchange, routingKey, false, false, pub); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
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
	r.connMu.Lock()
	if r.conn == nil || r.conn.IsClosed() {
		r.connMu.Unlock()
		return fmt.Errorf("connection is closed")
	}
	ch, err := r.conn.Channel()
	r.connMu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to open consumer channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
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
		if err := ch.QueueBind(q.Name, p, exchange, false, nil); err != nil {
			return fmt.Errorf("failed to bind queue to %s with pattern %s: %w", exchange, p, err)
		}
	}

	msgs, err := ch.ConsumeWithContext(ctx, q.Name, "", false, false, false, false, nil)
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
