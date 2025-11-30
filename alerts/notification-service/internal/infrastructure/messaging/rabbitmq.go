package messaging

import (
	"context"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
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
	if err := declareExchanges(firstCh, rabbitMQConfig.ExchangeType, rabbitMQConfig.Exchanges...); err != nil {
		return nil, err
	}
	broker.chPool <- firstCh

	log.Println("Successfully connected to RabittMQ")

	return broker, nil
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
		r.config.NotificationQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := ch.Qos(r.config.PrefetchCount, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	for _, p := range r.dispatcher.AvailablePatterns() {
		if err := ch.QueueBind(q.Name, p, exchange, false, nil); err != nil {
			return fmt.Errorf("failed to bind queue to %s with pattern %s: %w", exchange, p, err)
		}
	}
	log.Printf("Successfully bound %s exchange to %s queue", exchange, q.Name)

	msgs, err := ch.ConsumeWithContext(ctx, q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	log.Printf("Started consumer on queue %s", q.Name)

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

			go func(msg amqp.Delivery) {
				var wrapper event.EventWrapper
				pattern, data, err := wrapper.DecodedEventWrapper(msg.Body)
				if err != nil {
					log.Printf("Failed to decode event")
					msg.Nack(false, false)
					return
				}

				if err := r.dispatcher.HandleEvent(context.Background(), pattern, data); err != nil {
					log.Printf("Failed to handle event: %v", err)
					msg.Nack(false, true)
					return
				}

				msg.Ack(false)
			}(msg)
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
