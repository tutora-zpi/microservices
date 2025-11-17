package config

const (
	REDIS_ADDR     string = "REDIS_ADDR"
	REDIS_PASSWORD string = "REDIS_PASSWORD"
	REDIS_DB       string = "REDIS_DB"

	MONGO_HOST                string = "MONGO_HOST"
	MONGO_PORT                string = "MONGO_PORT"
	MONGO_USER                string = "MONGO_USER"
	MONGO_PASS                string = "MONGO_PASS"
	MONGO_URI                 string = "MONGO_URI"
	MONGO_DATABASE            string = "MONGO_DATABASE"
	MONGO_METADATA_COLLECTION string = "MONGO_METADATA_COLLECTION"

	RABBITMQ_DEFAULT_USER string = "RABBITMQ_DEFAULT_USER"
	RABBITMQ_DEFAULT_PASS string = "RABBITMQ_DEFAULT_PASS"
	RABBITMQ_HOST         string = "RABBITMQ_HOST"
	RABBITMQ_PORT         string = "RABBITMQ_PORT"
	RABBITMQ_URL          string = "RABBITMQ_URL"
	CHAT_EXCHANGE         string = "CHAT_EXCHANGE"
	MEETING_EXCHANGE      string = "MEETING_EXCHANGE"
	RECORDER_QUEUE        string = "RECORDER_QUEUE"

	APP_PORT string = "APP_PORT"
	APP_ENV  string = "APP_ENV"

	JWKS_URL string = "JWKS_URL"

	WS_GATEWAY_URL string = "WS_GATEWAY_URL"

	AWS_ACCESS_KEY_ID     string = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY string = "AWS_SECRET_ACCESS_KEY"
	AWS_BUCKET_NAME       string = "AWS_BUCKET_NAME"
	AWS_REGION            string = "AWS_REGION"

	CLIENT_SECRET string = "CLIENT_SECRET"
)
