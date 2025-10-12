package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	eventhandler "notification-serivce/internal/app/event_handler"
	"notification-serivce/internal/app/service"
	"notification-serivce/internal/config"
	classinvitation "notification-serivce/internal/domain/event/class_invitation"
	meetinginvitation "notification-serivce/internal/domain/event/meeting_invitation"
	"notification-serivce/internal/infrastructure/bus"
	"notification-serivce/internal/infrastructure/messaging"
	"notification-serivce/internal/infrastructure/mongo"
	notificationmanager "notification-serivce/internal/infrastructure/notification_manager"
	handlers "notification-serivce/internal/infrastructure/rest/v1"
	"notification-serivce/internal/infrastructure/security"
	"notification-serivce/internal/infrastructure/server"

	"github.com/joho/godotenv"
)

var rabbitmqConfig messaging.RabbitConfig
var mongoConfig mongo.MongoConfig

func setupMongoConfig() {
	mongoConfig = mongo.MongoConfig{
		User:       os.Getenv(config.MONGO_USER),
		Pass:       os.Getenv(config.MONGO_PASS),
		Host:       os.Getenv(config.MONGO_HOST),
		Port:       os.Getenv(config.MONGO_PORT),
		Uri:        os.Getenv(config.MONGO_URI),
		DbName:     os.Getenv(config.MONGO_DB_NAME),
		Collection: os.Getenv(config.MONGO_COLLECTION),
	}
}

func setupRabbitMQConfig() {
	rabbitmqConfig = messaging.RabbitConfig{
		User:                 os.Getenv(config.RABBITMQ_DEFAULT_USER),
		Pass:                 os.Getenv(config.RABBITMQ_DEFAULT_PASS),
		Host:                 os.Getenv(config.RABBITMQ_HOST),
		Port:                 os.Getenv(config.RABBITMQ_PORT),
		URL:                  os.Getenv(config.RABBITMQ_URL),
		Retries:              3,
		NotificationExchange: os.Getenv(config.NOTIFICATION_EXCHANGE),
	}
}

func init() {
	env := os.Getenv(config.APP_ENV)
	if env == "" || env == "localhost" || env == "127.0.0.1" {
		_ = godotenv.Load(".env.local")
	}

	setupRabbitMQConfig()
	setupMongoConfig()

	security.FetchSignKey()
}

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dispatcher := bus.NewDispatcher()

	manager := notificationmanager.NewManager()
	manager.EnableBuffering(1000, 30*time.Minute)

	repo, err := mongo.NewNotificationRepository(rootCtx, mongoConfig)
	if err != nil {
		log.Panicf("failed to create repository: %v", err)
	}
	defer func() {
		if err := repo.Close(rootCtx); err != nil {
			log.Panicf("failed to close repository: %v", err)
		}
	}()

	service := service.NewNotificationSerivce(repo)

	broker, err := messaging.NewRabbitBroker(rabbitmqConfig, dispatcher)
	if err != nil {
		log.Panic(err)
	}

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

	httpServer := server.NewServer(handlers.NewRouter(manager, service))

	go func() {
		if err = broker.Consume(rootCtx); err != nil {
			log.Println(err)
		}
	}()

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
