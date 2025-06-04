package main

import (
	"os"
	"voice-service/internal/domain/model"
	"voice-service/internal/infrastructure/config"
	"voice-service/internal/infrastructure/database"
	"voice-service/internal/infrastructure/messaging"
	repoimpl "voice-service/internal/infrastructure/repo-impl"
	"voice-service/internal/infrastructure/rest"
	"voice-service/internal/infrastructure/ws"

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

	broker := messaging.NewRabbitBroker(config.NewRabbitConfig(os.Getenv(config.RABBITMQ_URL), "meeting", 4))
	defer broker.Close()

	gw := ws.NewGateway(os.Getenv(config.JWT_SECRET))
	repo := repoimpl.NewVoiceMeetingRepository(postgres)

	incjectable := config.Incjectable{
		Broker:   broker,
		Recorder: nil,
		Repo:     repo,
		Gateway:  gw,
	}

	router := rest.NewRouter(&incjectable)
	server := rest.NewServer(router)

	server.StartAndListen()

	// to inject
	// recorder := recorder.NewRecorder(nil)

}
