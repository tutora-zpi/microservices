package main

import (
	"context"
	"log"
	eventhandler "notification-serivce/internal/app/event_handler"
	"notification-serivce/internal/app/service"
	"notification-serivce/internal/config"
	classinvitation "notification-serivce/internal/domain/event/class_invitation"
	meetinginvitation "notification-serivce/internal/domain/event/meeting_invitation"
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
	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		if err := godotenv.Load(); err != nil {
			log.Panic(".env* file not found. Please check path or provide one.")
		}
	}

	// security.FetchSignKey()
}

func main() {
	dispatcher := bus.NewDispatcher()
	database := database.Connect()
	defer database.Close()

	manager := notificationmanager.NewManager()
	manager.EnableBuffering(1000, 30*time.Minute)
	repo := repository.NewNotificationRepository(database)
	service := service.NewNotificationSerivce(repo)
	broker := messaging.NewRabbitBroker(dispatcher)
	defer broker.Close()

	dispatcher.Register(
		&classinvitation.ClassInvitationCreatedEvent{},
		eventhandler.NewClassInvitationCreatedHandler(broker, repo),
	)

	dispatcher.Register(
		&classinvitation.ClassInvitationReadyEvent{},
		eventhandler.NewClassInvitationReadyHandler(manager, repo),
	)

	dispatcher.Register(
		&classinvitation.UserDetailsRespondedEvent{},
		eventhandler.NewUserDetailsResponsedHandler(broker, repo),
	)

	dispatcher.Register(
		&meetinginvitation.MeetingStartedEvent{},
		eventhandler.NewMeetingInvitationReadyEventHandler(manager, repo),
	)

	server := server.NewServer(handlers.NewRouter(manager, service))

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
