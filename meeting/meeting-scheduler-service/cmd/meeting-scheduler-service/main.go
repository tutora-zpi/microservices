package main

import (
	"context"
	"log"
	"meeting-scheduler-service/internal/app/usecase"
	"meeting-scheduler-service/internal/infrastructure/config"
	"meeting-scheduler-service/internal/infrastructure/messaging"
	"meeting-scheduler-service/internal/infrastructure/redis"
	"meeting-scheduler-service/internal/infrastructure/rest"
	"meeting-scheduler-service/internal/infrastructure/rest/v1/handlers"
	"meeting-scheduler-service/internal/infrastructure/security"
	"meeting-scheduler-service/internal/infrastructure/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		if err := godotenv.Load(); err != nil {
			log.Panic(".env* file not found. Please check path or provide one.")
		}
	}

	security.FetchSignKey()
}

func main() {
	broker := messaging.NewRabbitBroker()
	defer broker.Close()
	repo := redis.NewMeetingRepo()
	defer repo.Close()

	meetingManager := usecase.NewMeetingManager(broker, repo)

	router := rest.NewRouter(handlers.NewManageMeetingHandler(meetingManager))

	server := server.NewServer(router)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.StartAndListen(); err != nil {
			log.Printf("server stopped with error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()

	if err := server.GracefulShutdown(context.Background()); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
