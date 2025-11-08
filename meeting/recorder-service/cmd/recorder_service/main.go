package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	eventhandlers "recorder-service/internal/app/event_handlers"
	factoryimpl "recorder-service/internal/app/factory_impl"
	"recorder-service/internal/app/interfaces"
	"recorder-service/internal/app/service"
	"recorder-service/internal/config"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/ws_event/general"
	"recorder-service/internal/domain/ws_event/rtc"
	"recorder-service/internal/infrastructure/bus"
	"recorder-service/internal/infrastructure/cache"
	"recorder-service/internal/infrastructure/messaging"
	"recorder-service/internal/infrastructure/mongo"
	"recorder-service/internal/infrastructure/redis"
	repoimpl "recorder-service/internal/infrastructure/repository"
	"recorder-service/internal/infrastructure/rest"
	"recorder-service/internal/infrastructure/rest/v1/handlers"
	"recorder-service/internal/infrastructure/s3"
	"recorder-service/internal/infrastructure/security"
	"recorder-service/internal/infrastructure/server"
	"recorder-service/internal/infrastructure/webrtc/writer"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	redisdb "github.com/redis/go-redis/v9"
	mongodb "go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		_ = godotenv.Load(".env.local")
	}

	security.FetchSignKey()
}

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	initCtx, cancel := context.WithTimeout(rootCtx, 3*time.Second)
	defer cancel()

	var errors chan error = make(chan error, 3)

	var mongoClient *mongodb.Client
	var redisClient *redisdb.Client

	var mongoConfig mongo.MongoConfig = *mongo.NewMongoConfig(time.Second * 5)
	var redisConfig redis.RedisConfig = *redis.NewRedisConfig(time.Second * 5)
	var rabbitMQConfig messaging.RabbitConfig = *messaging.NewRabbitMQConfig(time.Second*4, 10)

	var s3Service s3.S3Service

	var dispatcher bus.Dispachable = bus.NewDispatcher()

	var broker interfaces.Broker

	var wg sync.WaitGroup

	wg.Go(func() {
		var err error
		s3Service, err = s3.NewS3Service(initCtx, os.Getenv(config.AWS_BUCKET_NAME))
		if err != nil {
			errors <- err
		}
	})

	wg.Go(func() {
		var err error
		broker, err = messaging.NewRabbitBroker(rabbitMQConfig, dispatcher)
		if err != nil {
			errors <- err
		}
	})

	wg.Go(func() {
		var err error
		mongoClient, err = mongo.NewMongoClient(initCtx, mongoConfig)

		if err != nil {
			errors <- err
		}
	})

	wg.Go(func() {
		var err error
		redisClient, err = redis.NewRedis(initCtx, redisConfig)

		if err != nil {
			errors <- err
		}
	})

	wg.Wait()

	close(errors)

	for err := range errors {
		log.Fatalf("Something went wrong with services: %v", err)
	}

	defer func() {
		if redisClient != nil {
			log.Println("Closing redis connection...")
			_ = redisClient.Close()
		}

		if mongoClient != nil {
			mongo.Close(rootCtx, mongoClient, time.Second*5)
		}
	}()

	voiceRepo := repoimpl.NewVoiceMeetingRepository(mongoClient, mongoConfig)
	botRepo := cache.NewBotRepository(redisClient)

	voiceSessionService := service.NewVoiceSessionService(voiceRepo)

	clientFactory := factoryimpl.NewClientFactory(dispatcher)
	recorderFactory := factoryimpl.NewRecorderFactory()

	botService := service.NewBotService(botRepo, recorderFactory, clientFactory)

	dispatcher.Register(&event.MeetingStartedEvent{}, eventhandlers.NewMeetingStartedHandler(voiceRepo))
	dispatcher.Register(&event.StopRecordingMeetingEvent{}, eventhandlers.NewStopRecordingMeetingHandler(botService, voiceRepo, s3Service, broker, rabbitMQConfig.MeetingExchange))
	dispatcher.Register(&event.RecordMeetingEvent{}, eventhandlers.NewRecorderMeetingHandler(botService, writer.NewLocalWriter))

	dispatcher.Register(&rtc.OfferWSEvent{}, eventhandlers.NewOfferHandler(botService, writer.NewLocalWriter))
	dispatcher.Register(&rtc.IceCandidateWSEvent{}, eventhandlers.NewIceCandidateHandler(botService))
	dispatcher.Register(&general.RoomUsersWSEvent{}, eventhandlers.NewRoomUsersHandler(botService, writer.NewLocalWriter))

	handlers := handlers.NewHandler(voiceSessionService)

	router := rest.NewRouter(handlers)
	server := server.NewServer(router)

	go func() {
		if err := broker.Consume(rootCtx, rabbitMQConfig.MeetingExchange); err != nil {
			log.Printf("Consuming error: %v", err)
		}
	}()

	go func() {
		if err := server.StartAndListen(); err != nil {
			log.Printf("server stopped with error: %v", err)
			stop()
		}
	}()

	<-rootCtx.Done()

	if err := server.GracefulShutdown(context.Background()); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
