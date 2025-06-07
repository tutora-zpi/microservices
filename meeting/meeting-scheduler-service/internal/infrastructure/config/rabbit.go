package config

type RabbitConfig struct {
	// Connection string to RabbitMQ server
	Connstr string

	// How many times to retry connection
	Retries int
}

func NewRabbitConfig(connstr string, retries int) RabbitConfig {
	if connstr == "" {
		panic("RabbitMQ connection string cannot be empty")
	}

	return RabbitConfig{
		Connstr: connstr,
		Retries: retries,
	}
}
