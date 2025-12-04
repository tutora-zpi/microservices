package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	eventhandler "notification-serivce/internal/app/event_handler"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/app/service"
	"notification-serivce/internal/config"
	classinvitation "notification-serivce/internal/domain/event/class_invitation"
	meetinginvitation "notification-serivce/internal/domain/event/meeting_invitation"
	"notification-serivce/internal/infrastructure/bus"
	"notification-serivce/internal/infrastructure/messaging"
	"notification-serivce/internal/infrastructure/mongo"
	notificationmanager "notification-serivce/internal/infrastructure/notification_manager"
	"notification-serivce/internal/infrastructure/repository"
	handlers "notification-serivce/internal/infrastructure/rest/v1"
	"notification-serivce/internal/infrastructure/security"
	"notification-serivce/internal/infrastructure/server"

	mongodb "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv(config.APP_ENV)
	if env == "" || env == "localhost" || env == "127.0.0.1" {
		_ = godotenv.Load(".env.local")
	}

	security.FetchSignKey()
}

func main() {
	var rabbitmqConfig messaging.RabbitConfig = *messaging.NewRabbitMQConfig(time.Second*5, 10)
	var mongoConfig mongo.MongoConfig = *mongo.NewMongoConfig()
	var broker interfaces.Broker
	var mongoClient *mongodb.Client
	var wg sync.WaitGroup
	var errors chan error = make(chan error, 2)

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dispatcher := bus.NewDispatcher()

	manager := notificationmanager.NewManager()
	manager.EnableBuffering(1000, 30*time.Minute)

	wg.Go(func() {
		var err error
		mongoClient, err = mongo.NewMongoClient(rootCtx, mongoConfig, time.Second*5)
		if err != nil {
			errors <- err
		}
	})

	wg.Go(func() {
		var err error
		broker, err = messaging.NewRabbitBroker(rabbitmqConfig, dispatcher)
		if err != nil {
			errors <- err
		}
	})

	wg.Wait()

	close(errors)

	for err := range errors {
		log.Fatalf("Error during initialization: %v", err)
	}

	cleanCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer func() {
		mongo.Close(cleanCtx, mongoClient, time.Second*5)

		cancel()
	}()

	repo := repository.NewNotificationRepository(mongoClient, mongoConfig)

	service := service.NewNotificationSerivce(repo)

	dispatcher.Register(
		&classinvitation.ClassInvitationCreatedEvent{},
		eventhandler.NewClassInvitationCreatedHandler(manager, repo),
	)
	dispatcher.Register(
		&meetinginvitation.MeetingStartedEvent{},
		eventhandler.NewMeetingInvitationReadyEventHandler(manager, repo),
	)
	dispatcher.Register(
		&meetinginvitation.PlannedMeetingEvent{},
		eventhandler.NewMeetingPlannedHandler(manager, repo),
	)
	dispatcher.Register(
		&meetinginvitation.MeetingEndedEvent{},
		eventhandler.NewMeetingEndedHandler(manager, repo),
	)
	dispatcher.Register(
		&classinvitation.ClassInvitationAcceptedEvent{},
		eventhandler.NewMeetingEndedHandler(manager, repo),
	)

	httpServer := server.NewServer(handlers.NewRouter(manager, service))

	for _, exchange := range rabbitmqConfig.Exchanges {
		go func() {
			if err := broker.Consume(rootCtx, exchange); err != nil {
				log.Println(err)
			}
		}()
	}

	go func() {
		if err := httpServer.StartAndListen(); err != nil {
			stop()
		}
	}()

	<-rootCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpServer.GracefulShutdown(shutdownCtx)

}
