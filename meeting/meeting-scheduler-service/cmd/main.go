package main

import (
	"meeting-scheduler-service/internal/app/usecase"
	"meeting-scheduler-service/internal/infrastructure/config"
	"meeting-scheduler-service/internal/infrastructure/handlers"
	"meeting-scheduler-service/internal/infrastructure/messaging"
	"meeting-scheduler-service/internal/infrastructure/rest"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	broker := messaging.NewRabbitBroker(config.NewRabbitConfig(os.Getenv(config.RABBITMQ_URL), 4))
	defer broker.Close()

	meetingManager := usecase.NewMeetingManager(broker)

	router := rest.NewRouter(handlers.NewManageMeetingHandler(meetingManager))

	server := rest.NewServer(router)

	server.StartAndListen()
}
