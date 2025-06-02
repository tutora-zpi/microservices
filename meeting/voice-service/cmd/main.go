package main

import (
	"log"
	"net/http"
	"os"
	"voice-service/internal/domain/model"
	"voice-service/internal/infrastructure/config"
	"voice-service/internal/infrastructure/database"
	"voice-service/internal/infrastructure/messaging"
	"voice-service/internal/infrastructure/ws"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	postgres := database.NewPostgres(config.NewPostgresConfig(os.Getenv("POSTGRES_URL"), 4, nil, model.VoiceSession{}))
	defer postgres.Close()

	broker := messaging.NewRabbitBroker(config.NewRabbitConfig(os.Getenv("RABBITMQ_URL"), "meeting", 4))
	defer broker.Close()

	gw := ws.NewGateway()

	// to inject
	// recorder := recorder.NewRecorder(nil)

	http.HandleFunc("/ws", gw.HandleWS)

	log.Println("Signaling server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
