package config

const (
	JWKS_URL string = "JWKS_URL"
	APP_ENV  string = "APP_ENV"
	APP_NAME string = "APP_NAME"
	APP_PORT string = "APP_PORT"

	// mongo creds
	MONGO_HOST       string = "MONGO_HOST"
	MONGO_PASS       string = "MONGO_PASS"
	MONGO_DB_NAME    string = "MONGO_DB_NAME"
	MONGO_COLLECTION string = "MONGO_COLLECTION"
	MONGO_USER       string = "MONGO_USER"
	MONGO_PORT       string = "MONGO_PORT"
	MONGO_URI        string = "MONGO_URI"

	//rabbitmq creds
	RABBITMQ_DEFAULT_USER string = "RABBITMQ_DEFAULT_USER"
	RABBITMQ_DEFAULT_PASS string = "RABBITMQ_DEFAULT_PASS"
	RABBITMQ_HOST         string = "RABBITMQ_HOST"
	RABBITMQ_PORT         string = "RABBITMQ_PORT"
	RABBITMQ_URL          string = "RABBITMQ_URL"

	EVENT_EXCHANGE_QUEUE_NAME string = "EVENT_EXCHANGE_QUEUE_NAME"
	FRONTEND_URL              string = "FRONTEND_URL"
)
