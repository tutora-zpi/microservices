package main

import (
	"os"
	"voice-service/internal/domain/model"
	"voice-service/internal/infrastructure/config"
	"voice-service/internal/infrastructure/database"
	"voice-service/internal/infrastructure/messaging"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	postgres := database.NewPostgres(config.NewPostgresConfig(os.Getenv(config.POSTGRES_URL), 4, nil, model.VoiceSession{}))
	defer postgres.Close()

	broker := messaging.NewRabbitBroker(config.NewRabbitConfig(os.Getenv(config.RABBITMQ_URL), 4))
	defer broker.Close()

	// gw := ws.NewGateway(os.Getenv(config.JWT_SECRET))
	// repo := repoimpl.NewVoiceMeetingRepository(postgres)

	// router := rest.NewRouter(gw)
	// server := rest.NewServer(router)

	// server.StartAndListen()

	// to inject
	// recorder, err := recorder.NewRecorder(nil, writer.NewLocalWriter)

}
