package config

type RabbitConfig struct {
	// Connection string to RabbitMQ server
	Connstr string

	// Name of queue to consume messages from or publish
	Queue string

	// How many times to retry connection
	Retries int
}

func NewRabbitConfig(connstr string, queue string, retries int) RabbitConfig {
	if connstr == "" {
		panic("RabbitMQ connection string cannot be empty")
	}
	if queue == "" {
		panic("RabbitMQ queue name cannot be empty")
	}

	return RabbitConfig{
		Connstr: connstr,
		Queue:   queue,
		Retries: retries,
	}
}
