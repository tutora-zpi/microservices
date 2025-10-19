package main

import (
	eventhandler "chat-service/internal/app/event_handler"
	"chat-service/internal/app/interfaces"
	"chat-service/internal/app/service"
	"chat-service/internal/config"
	"chat-service/internal/domain/event"
	"chat-service/internal/infrastructure/bus"
	"chat-service/internal/infrastructure/messaging"
	mongoConn "chat-service/internal/infrastructure/mongo"
	"chat-service/internal/infrastructure/repository"
	"chat-service/internal/infrastructure/rest"
	"chat-service/internal/infrastructure/rest/v1/handlers"
	"chat-service/internal/infrastructure/server"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatal(".env* file not found. Please check path or provide one.")
		}
	}

	// security.FetchSignKey()
}

func main() {
	var mongoConfig *mongoConn.MongoConfig = mongoConn.NewMongoConfig()
	var rabbitMQConfig *messaging.RabbitConfig = messaging.NewRabbitConfig()

	var mongoClient *mongo.Client
	var dispatcher bus.Dispachable = bus.NewDispatcher()
	var broker interfaces.Broker

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	var errors chan error = make(chan error, 2)

	wg.Go(func() {
		var err error
		mongoClient, err = mongoConn.NewMongoClient(rootCtx, mongoConfig.URL(), time.Second*5)
		if err != nil {
			errors <- err
		}
	})

	wg.Go(func() {
		var err error
		broker, err = messaging.NewRabbitBroker(*rabbitMQConfig, dispatcher)
		if err != nil {
			errors <- err
		}
	})

	wg.Wait()

	close(errors)

	for err := range errors {
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	chatRepo := repository.NewChatRepository(mongoClient, *mongoConfig)
	messageRepo := repository.NewMessageRepository(mongoClient, *mongoConfig)

	dispatcher.Register(&event.SendMessageEvent{}, eventhandler.NewSendMessageHandler(messageRepo))
	dispatcher.Register(&event.MeetingStartedEvent{}, eventhandler.NewMeetingStartedHandler(chatRepo))
	dispatcher.Register(&event.ReactMessageOnEvent{}, eventhandler.NewReactHandler(messageRepo))
	dispatcher.Register(&event.ReplyOnMessageEvent{}, eventhandler.NewReplyHandler(messageRepo))

	log.Println("All services initialized successfully")

	go func(ctx context.Context) {
		if err := broker.Consume(ctx, rabbitMQConfig.ChatExchange); err != nil {
			log.Printf("Consumer error: %v", err)

		}
	}(rootCtx)

	go func(ctx context.Context) {
		if err := broker.Consume(ctx, rabbitMQConfig.MeetingExchange); err != nil {
			log.Printf("Consumer error: %v", err)

		}
	}(rootCtx)

	chatService := service.NewChatService(chatRepo)
	messageService := service.NewMessageService(messageRepo)

	handlers := handlers.NewHandlers(chatService, messageService)

	router := rest.NewRouter(handlers)
	server := server.NewServer(router)
	go func() {
		if err := server.StartAndListen(); err != nil {
			log.Printf("server stopped with error: %v", err)
			stop()
		}
	}()

	<-rootCtx.Done()
	log.Println("Shutting down...")

	cleanCtx, cancel := context.WithTimeout(rootCtx, 5*time.Second)
	defer cancel()

	if err := server.GracefulShutdown(cleanCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	mongoConn.Close(cleanCtx, mongoClient, time.Second*1)

}
