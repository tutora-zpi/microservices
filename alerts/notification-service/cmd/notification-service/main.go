package main

import (
	"context"
	"log"
	"notification-serivce/internal/app/usecase"
	"notification-serivce/internal/domain/event"
	"notification-serivce/internal/domain/query"
	"notification-serivce/internal/infrastructure/bus"
	"notification-serivce/internal/infrastructure/database"
	"notification-serivce/internal/infrastructure/messaging"
	notificationmanager "notification-serivce/internal/infrastructure/notification_manager"
	"notification-serivce/internal/infrastructure/repository"
	handlers "notification-serivce/internal/infrastructure/rest/v1"
	"notification-serivce/internal/infrastructure/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Panic(".env* file not found. Please check path or provide one.")
	}

	// security.FetchSignKey()
}

func main() {
	dispatcher := bus.NewDispatcher()
	queryBus := bus.NewQueryBus()
	database := database.Connect()
	defer database.Close()

	notificationManager := notificationmanager.NewManager()
	notificationManager.EnableBuffering(1000, 10*time.Minute)
	repo := repository.NewNotificationRepository(database)

	queryBus.Register(
		&query.FetchNotificationsQuery{},
		usecase.NewFetchNotificationsHandler(repo),
	)

	dispatcher.Register(
		&event.ClassInvitationEvent{},
		usecase.NewClassInvitationHandler(notificationManager, repo),
	)

	server := server.NewServer(handlers.NewRouter(notificationManager, queryBus))

	broker := messaging.NewRabbitBroker(dispatcher)
	defer broker.Close()

	go broker.Consume()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// HTTP

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
